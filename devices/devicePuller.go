package devices

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/altucor/http-knocker/deviceCommand"
	"github.com/altucor/http-knocker/deviceCommon"
	"github.com/altucor/http-knocker/deviceResponse"
	"github.com/altucor/http-knocker/firewallCommon"
	"github.com/altucor/http-knocker/logging"

	"github.com/gorilla/mux"
)

type vfwc uint16

const (
	VFWC_STATE_PENDING_ADD    vfwc = 1
	VFWC_STATE_ADDED          vfwc = 2
	VFWC_STATE_PENDING_REMOVE vfwc = 3
)

type virtualFirewallCmd struct {
	id    string
	cmd   deviceCommon.IDeviceCommand
	state vfwc
}

func (ctx virtualFirewallCmd) toMap() map[string]interface{} {
	vfcmd := make(map[string]interface{})
	vfcmd["id"] = ctx.id
	vfcmd["command"] = ctx.cmd.ToMap()
	return vfcmd
}

type virtualFirewall struct {
	mu    sync.Mutex
	cmds  []virtualFirewallCmd
	rules []firewallCommon.FirewallRule
}

/*
Important:
Instead of storing rules in to virtual firewall.
Save actual commands in virtual firewall state.
Because communication with remote client should be with commands, not rules
Remote client should know how to interpret commands in rules
and how to execute commands for custom firewall
*/

func generateCmdId() string {
	h := sha1.New()
	h.Write([]byte(time.Now().String()))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (ctx *virtualFirewall) Add(cmd deviceCommand.Add) (deviceResponse.Add, error) {
	// Should add new commands to pending list until they will be accepted by remote firewall
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.cmds = append(ctx.cmds, virtualFirewallCmd{
		id:    generateCmdId(),
		cmd:   cmd,
		state: VFWC_STATE_PENDING_ADD,
	})
	return deviceResponse.Add{}, nil
}

func (ctx *virtualFirewall) Get(cmd deviceCommand.Get) (deviceResponse.Get, error) {
	// Should return with list of accepted virtual firewall rules
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	return deviceResponse.GetFromRuleList(ctx.rules)
}

func (ctx *virtualFirewall) Remove(cmd deviceCommand.Remove) (deviceResponse.Remove, error) {
	// Should mark rules from accepted list as pending for removal, but not remove them
	// Only really remove them when remote firewall will approve this
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.cmds = append(ctx.cmds, virtualFirewallCmd{
		id:    generateCmdId(),
		cmd:   cmd,
		state: VFWC_STATE_PENDING_ADD,
	})
	return deviceResponse.Remove{}, nil
}

func (ctx *virtualFirewall) getLastUpdates(count uint64) (string, error) {
	// Here we respond only with pending changes for remote firewall
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	var cmds []map[string]interface{}
	for _, item := range ctx.cmds {
		if count != 0 && uint64(len(cmds)) >= count {
			break
		}
		cmds = append(cmds, item.toMap())
	}
	jsonBytes, err := json.Marshal(cmds)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func (ctx *virtualFirewall) acceptUpdates(acceptedRules []string) error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	for _, acceptedRule := range acceptedRules {
		for i, item := range ctx.cmds {
			if item.id == acceptedRule {
				// TODO: Here remove cmd and break cycle
				ctx.cmds = append(ctx.cmds[:i], ctx.cmds[i+1:]...)
				break
			}
		}
	}
	return nil
}

func (ctx *virtualFirewall) pushRuleSet(rules []firewallCommon.FirewallRule) error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	// TODO: Here overwrite rule set with rules from remote firewall
	ctx.rules = rules
	return nil
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
	if err := r.ParseForm(); err != nil {
		logging.CommonLog().Debugf("ParseForm() err: %v", err)
		return
	}
	acceptedRulesJson := r.FormValue("accepted_rules")
	logging.CommonLog().Debugf("Accepted rules: %s", acceptedRulesJson)
	var acceptedRules []string
	err := json.Unmarshal([]byte(r.FormValue("accepted_rules")), &acceptedRules)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logging.CommonLog().Debugf("Accepted rules: %v", acceptedRules)
	err = ctx.firewallState.acceptUpdates(acceptedRules)
	if err != nil {
		logging.CommonLog().Error(err)
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "500\n")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (ctx *DevicePuller) pushRulesSet(w http.ResponseWriter, r *http.Request) {
	logging.CommonLog().Info("[devicePuller] called pushRulesSet")
	if err := r.ParseForm(); err != nil {
		logging.CommonLog().Debugf("ParseForm() err: %v", err)
		return
	}
	frwRulesJson := r.FormValue("rules")
	logging.CommonLog().Debugf("Frw rules: %s", frwRulesJson)
	var frwRules []firewallCommon.FirewallRule
	err := json.Unmarshal([]byte(r.FormValue("rules")), &frwRules)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logging.CommonLog().Debugf("Frw rules: %s", frwRules)
	err = ctx.firewallState.pushRuleSet(frwRules)
	if err != nil {
		logging.CommonLog().Error(err)
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "500\n")
		return
	}
	w.WriteHeader(http.StatusOK)

	return
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
	ctx.router.NotFoundHandler = http.HandlerFunc(http_not_found_handler)
	pullerRouter := ctx.router.PathPrefix(ctx.config.Endpoint).Subrouter()
	pullerRouter.HandleFunc("/getLastUpdates", ctx.getLastUpdates).Methods("GET")
	pullerRouter.HandleFunc("/acceptUpdates", ctx.acceptUpdates).Methods("POST")
	pullerRouter.HandleFunc("/pushRulesSet", ctx.pushRulesSet).Methods("POST")
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
		return ctx.firewallState.Add(command.(deviceCommand.Add))
	case deviceCommon.DeviceCommandGet:
		return ctx.firewallState.Get(command.(deviceCommand.Get))
	case deviceCommon.DeviceCommandRemove:
		return ctx.firewallState.Remove(command.(deviceCommand.Remove))
	}

	return nil, errors.New("invalid command type")
}
