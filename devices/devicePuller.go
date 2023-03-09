package devices

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"net/http"
	"strconv"
	"time"

	"github.com/altucor/http-knocker/deviceCommand"
	"github.com/altucor/http-knocker/deviceCommon"
	"github.com/altucor/http-knocker/firewallCommon"
	"github.com/altucor/http-knocker/logging"

	"github.com/gorilla/mux"
)

type vfrs uint16

const (
	VFR_STATE_PENDING_ADD    vfrs = 1
	VFR_STATE_ADDED          vfrs = 2
	VFR_STATE_PENDING_REMOVE vfrs = 3
)

type virtualFirewallRule struct {
	rule  firewallCommon.FirewallRule
	state vfrs
}

type virtualFirewall struct {
	rules []virtualFirewallRule
}

func (ctx *virtualFirewall) Add(cmd deviceCommand.Add) error {
	// Should add new commands to pending list until they will be accepted by remote firewall
	rule := cmd.GetRule()
	rule.Id.SetValue(uint64(crc32.ChecksumIEEE([]byte(rule.Comment.GetString()))))
	ctx.rules = append(ctx.rules, virtualFirewallRule{
		rule:  rule,
		state: VFR_STATE_PENDING_ADD,
	})
	return nil
}

func (ctx *virtualFirewall) Get(cmd deviceCommand.Get) ([]firewallCommon.FirewallRule, error) {
	// Should return with accepted clients by remote firewall
	var rules []firewallCommon.FirewallRule
	for _, item := range ctx.rules {
		if item.state == VFR_STATE_ADDED {
			rules = append(rules, item.rule)
		}
	}
	return rules, nil
}

func (ctx *virtualFirewall) Remove(cmd deviceCommand.Remove) error {
	// Should mark rules from accepted list as pending for removal, but not remove them
	// Only really remove them when remote firewall will approve this
	id := cmd.GetId()
	for iter, item := range ctx.rules {
		if item.rule.Id.GetValue() == id {
			switch item.state {
			case VFR_STATE_ADDED:
				ctx.rules[iter].state = VFR_STATE_PENDING_REMOVE
			case VFR_STATE_PENDING_ADD:
				ctx.rules = append(ctx.rules[:iter], ctx.rules[iter+1:]...)
				return nil
			case VFR_STATE_PENDING_REMOVE:
				// do nothing
			default:
				logging.CommonLog().Debug("Default case called what to do?")
			}
		}
	}
	return nil
}

func (ctx *virtualFirewall) getLastUpdates(count uint64) (string, error) {
	// Here we respond only with pending changes for remote firewall
	var rules []string
	var counter uint64 = 0
	for _, item := range ctx.rules {
		if counter >= count {
			break
		}
		switch item.state {
		case VFR_STATE_PENDING_ADD:
		case VFR_STATE_PENDING_REMOVE:
			// here convert each rule to json
			counter += 1
			jsoned, err := item.rule.ToJson()
			if err != nil {
				return "", err
			} else {
				rules = append(rules, jsoned)
			}
		}
	}

	jsonBytes, err := json.Marshal(rules)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

/*
Interface for web side:
- GET - getLastUpdates
	(optional "count" arg - how much entries to show)
	give json with all pending changes
	like Add Get Reomve commands with unique id's
- POST - acceptUpdates
	client notifies us about accepted/executed rules
	client packet should have id's of successfully executed rules
	after receiving this list virtual firewall should move add command to added clients
- POST - resetState
	Reset state for cases when we are in unsync with remote firewall
- POST - pushInitialState
	When httpKnocker becomes alive and remote
	firewall have some rules from previous run.
	It will allow to push initial state just to be in sync


Structure of commands and responses:
Generally all in JSON

Here is potential problems with removing some rules on remote side.
By design we don't need to know internals of remote firewall.
And here problem that we dont know under which ID
	actual client rules is added on remote firewall.
Because of this we need to introduce our own identification of rules.
Each rule after adding shoud get unique ID which also can be calculated
by remote firewall client.

As best for usage, open and non-compromised data we can use
hash of comment parts.

c = comment
uniqueRuleId = MD5(c.prefix + c.firewallName + c.timestamp + c.endpointHash)







Response for getLastUpdates
{
	"commands": [
		{
			"id": 425462
			"type": "add"
			"data": { here info about ip port chain action... }
		},
		{
			"id": 83452
			"type": "remove"
			"command": { here probably id of rule from which list? our or remote? }
		}
	]
}




*/

type ConnectionPuller struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Port     uint16 `yaml:"port"`
	Endpoint string `yaml:"endpoint"`
}

type DevicePuller struct {
	config             ConnectionPuller
	supportedProtocols []DeviceProtocol
	server             *http.Server
	router             *mux.Router
	firewallState      virtualFirewall
}

func (ctx *DevicePuller) getLastUpdates(w http.ResponseWriter, r *http.Request) {
	logging.CommonLog().Info("[devicePuller] called getLastUpdates")
	var itemsCount uint64 = 0
	if r.URL.Query().Get("count") != "" {
		count, err := strconv.ParseUint(r.URL.Query().Get("count"), 10, 64)
		if err != nil {
			logging.CommonLog().Info("[devicePuller] error decoding count getLastUpdates")
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "500\n")
			return
		} else {
			itemsCount = count
		}
	}
	rules, err := ctx.firewallState.getLastUpdates(itemsCount)
	if err != nil {
		logging.CommonLog().Info("[devicePuller] error preparing getLastUpdates")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", rules)
}

func (ctx *DevicePuller) acceptUpdates(w http.ResponseWriter, r *http.Request) {
	logging.CommonLog().Info("[devicePuller] called acceptUpdates")
	// Here mark rules with accepted statuses
}

func http_not_found_handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404\n")
}

func DevicePullerNew(cfg DeviceConnectionDesc) *DevicePuller {
	ctx := &DevicePuller{
		config: ConnectionPuller{
			Username: cfg.Username,
			Password: cfg.Password,
			Port:     cfg.Port,
			Endpoint: cfg.Endpoint,
		},
		supportedProtocols: []DeviceProtocol{
			PROTOCOL_ANY,
		},
		server: &http.Server{
			Addr:         "0.0.0.0" + ":" + fmt.Sprint(cfg.Port),
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		router: mux.NewRouter(),
	}
	ctx.router.HandleFunc("/", http_not_found_handler)
	ctx.router.HandleFunc(ctx.config.Endpoint+"getLastUpdates", ctx.getLastUpdates)
	ctx.router.HandleFunc(ctx.config.Endpoint+"acceptUpdates", ctx.acceptUpdates)
	return ctx
}

func (ctx *DevicePuller) GetSupportedProtocols() []DeviceProtocol {
	return ctx.supportedProtocols
}

func (ctx *DevicePuller) GetType() DeviceType {
	return DeviceTypePuller
}

func (ctx *DevicePuller) Start() error {
	ctx.server.Handler = ctx.router
	go func() {
		if err := ctx.server.ListenAndServe(); err != nil {
			logging.CommonLog().Error(err)
		}
	}()
	return nil
}

func (ctx *DevicePuller) Stop() error {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	ctx.server.Shutdown(ctxTimeout)
	return nil
}

func (ctx *DevicePuller) RunCommandWithReply(
	command deviceCommon.IDeviceCommand,
	proto DeviceProtocol,
) (deviceCommon.IDeviceResponse, error) {

	switch command.GetType() {
	case deviceCommon.DeviceCommandAdd:
		ctx.firewallState.Add(command.(deviceCommand.Add))
	case deviceCommon.DeviceCommandGet:
		ctx.firewallState.Get(command.(deviceCommand.Get))
	case deviceCommon.DeviceCommandRemove:
		ctx.firewallState.Remove(command.(deviceCommand.Remove))
	}

	logging.CommonLog().Error("Not implemented")
	return nil, errors.New("not Implemented")
}
