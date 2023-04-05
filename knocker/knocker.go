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
	Controllers map[string]*firewallControllers.InterfaceWrapper `yaml:"controllers"`
	Endpoints   map[string]*endpoint.Endpoint                    `yaml:"endpoints"`
	Knocks      map[string]*Knock                                `yaml:"knocks"`
}

func KnockerNewFromConfig(path string) (*Knocker, error) {
	knocker := &Knocker{
		WebServer:   nil,
		Devices:     make(map[string]devices.InterfaceWrapper),
		Controllers: make(map[string]*firewallControllers.InterfaceWrapper),
		Endpoints:   make(map[string]*endpoint.Endpoint),
		Knocks:      make(map[string]*Knock),
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		logging.CommonLog().Error("[Config] Error reading file: %s\n", err)
		return knocker, err
	}
	err = yaml.Unmarshal(bytes, &knocker)
	if err != nil {
		logging.CommonLog().Error("[Config] Error unmarshaling yaml file: %s\n", err)
		return knocker, err
	}

	for _, element := range knocker.Endpoints {
		element.SetDefaults()
	}

	// Setting Endpoint and Device to Controller
	for _, element := range knocker.Controllers {
		element.Controller.SetDevice(knocker.Devices[element.Controller.GetDeviceName()].Device)
		element.Controller.SetEndpoint(knocker.Endpoints[element.Controller.GetEndpointName()])
	}

	// Setting Controller and Endpoint to the Knock
	for _, element := range knocker.Knocks {
		element.SetController(knocker.Controllers[element.GetControllerName()].Controller)
		element.SetEndpoint(knocker.Endpoints[element.GetEndpointName()])
	}

	// Registering endpoints in webserver
	for _, element := range knocker.Knocks {
		knocker.WebServer.AddEndpoint(
			element.GetEndpoint().RegisterWithMiddlewares(
				element.GetHttpCallback,
			),
		)
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
