package knocker

import (
	"io/ioutil"
	"os"
	"os/signal"

	"gopkg.in/yaml.v3"

	"github.com/altucor/http-knocker/device/command"
	"github.com/altucor/http-knocker/device/response"
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
}

func KnockerNewFromConfig(path string) (*Knocker, error) {
	knocker := &Knocker{
		WebServer:   nil,
		Devices:     make(map[string]devices.InterfaceWrapper),
		Endpoints:   make(map[string]*endpoint.Endpoint),
		Controllers: make(map[string]*firewallControllers.InterfaceWrapper),
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		logging.CommonLog().Errorf("[Knocker] Error reading file: %s", err)
		return knocker, err
	}
	err = yaml.Unmarshal(bytes, &knocker)
	if err != nil {
		logging.CommonLog().Errorf("[Knocker] Error unmarshaling yaml file: %s", err)
		return knocker, err
	}

	controllerUrls := make(map[string]bool)
	for _, element := range knocker.Controllers {
		_, exist := controllerUrls[element.Config.Url]
		if exist {
			logging.CommonLog().Fatalf("[Knocker] Error detected several controllers with same URL: %s", element.Config.Url)
		}
		controllerUrls[element.Config.Url] = true
	}

	for _, element := range knocker.Endpoints {
		element.SetDefaults()
	}

	// Setting Endpoint and Device to Controller
	for _, element := range knocker.Controllers {
		element.Controller.SetDevice(knocker.Devices[element.Config.Device].Device)
		element.Controller.SetEndpoint(knocker.Endpoints[element.Config.Endpoint])
	}

	// Registering endpoints in webserver
	for _, element := range knocker.Controllers {
		knocker.WebServer.AddEndpoint(element.Controller.GetHttpCallback())
	}

	return knocker, nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func (ctx *Knocker) CleanupTrashRules() {
	knownIdentifiers := make([]string, 0)
	for _, element := range ctx.Controllers {
		knownIdentifiers = append(knownIdentifiers, element.Controller.GetHash())
	}

	for _, dev := range ctx.Devices {
		deviceResponse, err := dev.Device.RunCommandWithReply(command.GetNew())
		if err != nil {
			logging.CommonLog().Error("[Knocker] Error running command 'get rules' on device")
		}
		rules := deviceResponse.(*response.Get).GetRules()
		for _, rule := range rules {
			commentData, err := firewallControllers.FirewallCommentNewFromString(rule.Comment.GetString(), "-")
			if err != nil {
				logging.CommonLog().Error("[Knocker] Error processing firewall rule comment", err)
				continue
			}
			if commentData.GetEndpointHash() == "" {
				continue
			}

			if !contains(knownIdentifiers, commentData.GetEndpointHash()) {
				_, err := dev.Device.RunCommandWithReply(command.RemoveNew(rule.Id.GetValue()))
				if err != nil {
					logging.CommonLog().Error("[Knocker] Error removing unrelated rule:", rule)
				}
			}
		}
	}

}

func (ctx *Knocker) Start() {
	logging.CommonLog().Info("[Knocker] Starting...")
	for _, item := range ctx.Devices {
		item.Device.Start()
	}
	for _, item := range ctx.Controllers {
		item.Controller.Start()
	}
	ctx.CleanupTrashRules()
	ctx.WebServer.Start()
	logging.CommonLog().Info("[Knocker] Starting... DONE")
}

func (ctx *Knocker) Wait() {
	logging.CommonLog().Info("[Knocker] Waiting...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	logging.CommonLog().Info("[Knocker] Waiting... FINISHED")
}

func (ctx *Knocker) Stop() {
	logging.CommonLog().Info("[Knocker] Stopping...")
	ctx.WebServer.Stop()
	for _, item := range ctx.Controllers {
		item.Controller.Stop()
	}
	for _, item := range ctx.Devices {
		item.Device.Stop()
	}
	logging.CommonLog().Info("[Knocker] Stopping... DONE")
}
