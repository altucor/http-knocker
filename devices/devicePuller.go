package devices

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/altucor/http-knocker/device"
	"github.com/altucor/http-knocker/device/command"
	"github.com/altucor/http-knocker/firewallCommon"
	"github.com/altucor/http-knocker/logging"
	"gopkg.in/yaml.v3"

	"github.com/gorilla/mux"
)

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
	config        ConnectionPuller
	server        *http.Server
	router        *mux.Router
	firewallState VirtualFirewall
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
	commands, err := ctx.firewallState.getLastPendingCommands(itemsCount)
	if err != nil {
		logging.CommonLog().Info("[devicePuller] error preparing getLastUpdates")
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	logging.CommonLog().Debugf("Commands for update %s", commands)
	fmt.Fprintf(w, "%s", commands)
}

func (ctx *DevicePuller) setErrorCode(w http.ResponseWriter, errCode int, err error) {
	logging.CommonLog().Error(err)
	w.WriteHeader(errCode)
	fmt.Fprintf(w, "%d\n", errCode)
}

func (ctx *DevicePuller) getDataFromRequest(w http.ResponseWriter,
	r *http.Request, formKey string, jsonStruct any) error {
	if err := r.ParseForm(); err != nil {
		logging.CommonLog().Debug("ParseForm() err: ", err)
		return err
	}
	formValueData := r.FormValue(formKey)
	logging.CommonLog().Debug("formValueData: ", formValueData)
	err := json.Unmarshal([]byte(formValueData), jsonStruct)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *DevicePuller) acceptUpdates(w http.ResponseWriter, r *http.Request) {
	logging.CommonLog().Info("[devicePuller] called acceptUpdates")
	// Here mark rules with accepted statuses
	var acceptedRules []string
	if err := ctx.getDataFromRequest(w, r, "rules", &acceptedRules); err != nil {
		ctx.setErrorCode(w, 400, err)
		return
	}
	logging.CommonLog().Debug("Accepted rules: ", acceptedRules)
	err := ctx.firewallState.processAcceptedCommands(acceptedRules)
	if err != nil {
		ctx.setErrorCode(w, 500, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (ctx *DevicePuller) pushRulesSet(w http.ResponseWriter, r *http.Request) {
	// TODO: After fixing firewall rule interface, check this parsing
	logging.CommonLog().Info("[devicePuller] called pushRulesSet")
	// var frwRules []firewallCommon.FirewallRule
	var frwJsonRules []map[string]string
	if err := ctx.getDataFromRequest(w, r, "rules", &frwJsonRules); err != nil {
		ctx.setErrorCode(w, 400, err)
		return
	}

	var frwRules []firewallCommon.FirewallRule
	for _, elem := range frwJsonRules {
		rule := firewallCommon.FirewallRule{}
		rule.FromMap(elem)
		frwRules = append(frwRules, rule)
	}

	logging.CommonLog().Debug("Frw rules: ", frwRules)
	err := ctx.firewallState.pushRuleSet(frwRules)
	if err != nil {
		ctx.setErrorCode(w, 500, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func http_not_found_handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404\n")
}

func DevicePullerNew(cfg ConnectionPuller) *DevicePuller {
	ctx := &DevicePuller{
		config: ConnectionPuller{
			Username: cfg.Username,
			Password: cfg.Password,
			Port:     cfg.Port,
			Endpoint: cfg.Endpoint,
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
	pullerRouter := ctx.router.PathPrefix("/" + ctx.config.Endpoint).Subrouter()
	pullerRouter.HandleFunc("/getLastUpdates", ctx.getLastUpdates).Methods("GET")
	pullerRouter.HandleFunc("/acceptUpdates", ctx.acceptUpdates).Methods("POST")
	pullerRouter.HandleFunc("/pushRulesSet", ctx.pushRulesSet).Methods("POST")
	return ctx
}

func DevicePullerNewFromYaml(value *yaml.Node, protocol IFirewallRestProtocol) (*DevicePuller, error) {
	var cfg struct {
		Conn ConnectionPuller `yaml:"connection"`
	}
	if err := value.Decode(&cfg); err != nil {
		return nil, err
	}
	return DevicePullerNew(cfg.Conn), nil
}

func (ctx *DevicePuller) Start() error {
	logging.CommonLog().Info("[devicePuller] Starting...")
	ctx.server.Handler = ctx.router
	go func() {
		if err := ctx.server.ListenAndServe(); err != nil {
			logging.CommonLog().Error(err)
		}
	}()
	logging.CommonLog().Info("[devicePuller] Starting... DONE")
	return nil
}

func (ctx *DevicePuller) Stop() error {
	logging.CommonLog().Info("[devicePuller] Stopping...")
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	ctx.server.Shutdown(ctxTimeout)
	logging.CommonLog().Info("[devicePuller] Stopping... DONE")
	return nil
}

func (ctx *DevicePuller) RunCommandWithReply(
	cmd device.IDeviceCommand,
) (device.IDeviceResponse, error) {

	switch cmd.GetType() {
	case device.DeviceCommandAdd:
		return ctx.firewallState.Add(cmd.(command.Add))
	case device.DeviceCommandGet:
		return ctx.firewallState.Get(cmd.(command.Get))
	case device.DeviceCommandRemove:
		return ctx.firewallState.Remove(cmd.(command.Remove))
	}

	return nil, errors.New("invalid command type")
}
