package devices

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/altucor/http-knocker/deviceCommand"
	"github.com/altucor/http-knocker/deviceCommon"
	"github.com/altucor/http-knocker/firewallCommon"
	"github.com/altucor/http-knocker/firewallCommon/firewallField"
	"github.com/altucor/http-knocker/logging"

	"github.com/gorilla/mux"
)

type virtualFirewall struct {
	lastRuleId firewallField.Number
	rules      []firewallCommon.FirewallRule
}

func (ctx *virtualFirewall) Add(cmd deviceCommand.Add) error {
	// Should add new commands to pending list until they will be accepted by remote firewall
	rule := cmd.GetRule()
	rule.Id.SetValue(ctx.lastRuleId.GetValue())
	ctx.rules = append(ctx.rules, rule)
	return nil
}

func (ctx *virtualFirewall) Get(cmd deviceCommand.Get) ([]firewallCommon.FirewallRule, error) {
	// Should return with accepted clients by remote firewall
	return ctx.rules, nil
}

func (ctx *virtualFirewall) Remove(cmd deviceCommand.Remove) error {
	// Should mark rules from accepted list as pending for removal, but not remove them
	// Only really remove them when remote firewall will approve this
	id := cmd.GetId()
	ctx.rules = append(ctx.rules[:id], ctx.rules[id+1:]...)
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

func (ctx *DevicePuller) defaultHandler(w http.ResponseWriter, r *http.Request) {
	logging.CommonLog().Info("[devicePuller] called defaultHandler")
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
	ctx.router.HandleFunc(ctx.config.Endpoint, ctx.defaultHandler)
	return ctx
}

func (ctx *DevicePuller) GetSupportedProtocols() []DeviceProtocol {
	return ctx.supportedProtocols
}

func (ctx *DevicePuller) GetType() DeviceType {
	return DeviceTypePuller
}

func http_not_found_handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "1337\n")
}

func (ctx *DevicePuller) Start() error {
	ctx.router.HandleFunc(
		"/",
		http_not_found_handler,
	)
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
