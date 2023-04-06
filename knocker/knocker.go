package knocker

import (
	"io/ioutil"
	"os"
	"os/signal"

	"gopkg.in/yaml.v3"

	"github.com/altucor/http-knocker/devices"
	"github.com/altucor/http-knocker/endpoint"
	"github.com/altucor/http-knocker/firewallControllers"
	"github.com/altucor/http-knocker/logging"
	"github.com/altucor/http-knocker/webserver"
)

type Knocker struct {
	WebServer   *webserver.WebServer                             `yaml:"server"`
	Devices     map[string]devices.InterfaceWrapper              `yaml:"devices"`
	Endpoints   map[string]*endpoint.Endpoint                    `yaml:"endpoints"`
	Controllers map[string]*firewallControllers.InterfaceWrapper `yaml:"controllers"`
	// Knocks      map[string]*Knock                                `yaml:"knocks"`
}

func KnockerNewFromConfig(path string) (*Knocker, error) {
	knocker := &Knocker{
		WebServer:   nil,
		Devices:     make(map[string]devices.InterfaceWrapper),
		Controllers: make(map[string]*firewallControllers.InterfaceWrapper),
		Endpoints:   make(map[string]*endpoint.Endpoint),
		// Knocks:      make(map[string]*Knock),
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		logging.CommonLog().Errorf("[Knocker] Error reading file: %s\n", err)
		return knocker, err
	}
	err = yaml.Unmarshal(bytes, &knocker)
	if err != nil {
		logging.CommonLog().Errorf("[Knocker] Error unmarshaling yaml file: %s\n", err)
		return knocker, err
	}

	for _, element := range knocker.Endpoints {
		element.SetDefaults()
	}

	// Check for duplications in endpoints
	// To avoid cases when 2 endpoints assigned to one one controller or device
	// Two endpoints with identical URL cannot be used anyway,
	// because we can't register 2 endpoints under 1 url in webserver router.
	// But there is open question about how to track cases when 2 endpoints
	// with same port and protocol can affect same device
	for _, first := range knocker.Endpoints {
		for _, second := range knocker.Endpoints {
			if first.IsEqual(second) {
				logging.CommonLog().Error("Error found endpoints with duplicated params")
			}
		}
	}

	// Setting Endpoint and Device to Controller
	for _, element := range knocker.Controllers {
		element.Controller.SetDevice(knocker.Devices[element.Device].Device)
		element.Controller.SetEndpoint(knocker.Endpoints[element.Endpoint])
	}

	// Registering endpoints in webserver
	for _, element := range knocker.Controllers {
		knocker.WebServer.AddEndpoint(element.Controller.GetHttpCallback())
	}

	return knocker, nil
}

func (ctx *Knocker) Start() {
	logging.CommonLog().Info("[knocker] Starting...")
	for _, item := range ctx.Devices {
		item.Device.Start()
	}
	for _, item := range ctx.Controllers {
		item.Controller.Start()
	}
	ctx.WebServer.Start()
	logging.CommonLog().Info("[knocker] Starting... DONE")
}

func (ctx *Knocker) Wait() {
	logging.CommonLog().Info("[knocker] Waiting...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	logging.CommonLog().Info("[knocker] Waiting... FINISHED")
}

func (ctx *Knocker) Stop() {
	logging.CommonLog().Info("[knocker] Stopping...")
	ctx.WebServer.Stop()
	for _, item := range ctx.Controllers {
		item.Controller.Stop()
	}
	for _, item := range ctx.Devices {
		item.Device.Stop()
	}
	logging.CommonLog().Info("[knocker] Stopping... DONE")
}
