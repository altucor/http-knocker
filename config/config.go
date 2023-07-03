package config

import (
	"io/ioutil"

	"github.com/altucor/http-knocker/logging"
	"github.com/altucor/http-knocker/webserver"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Server webserver.ServerConfig
	// Devices   map[string]*devices.DeviceWrapper
	// Firewalls map[string]*firewallControllers.FirewallCfg
	// Endpoints map[string]*endpoint.EndpointCfg
	// Knocks    map[string]*knocker.KnockCfg
}

func (ctx *Configuration) SetDefaults() {
	// ctx.Server.SetDefaults()
	// for _, item := range ctx.Devices {
	// 	item.SetDefaults()
	// }
	// for _, item := range ctx.Firewalls {
	// 	item.SetDefaults()
	// }
	// for _, item := range ctx.Endpoints {
	// 	item.SetDefaults()
	// }
	// for _, item := range ctx.Knocks {
	// 	item.SetDefaults()
	// }
}

func (ctx *Configuration) Validate() error {
	// if err := ctx.Server.Validate(); err != nil {
	// 	return err
	// }
	// for _, item := range ctx.Devices {
	// 	if err := item.Validate(); err != nil {
	// 		return err
	// 	}
	// }
	// for _, item := range ctx.Firewalls {
	// 	if err := item.Validate(); err != nil {
	// 		return err
	// 	}
	// }
	// for _, item := range ctx.Endpoints {
	// 	if err := item.Validate(); err != nil {
	// 		return err
	// 	}
	// }
	// for _, item := range ctx.Knocks {
	// 	if err := item.Validate(); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

func ConfigurationNew(path string) (Configuration, error) {
	cfg := Configuration{}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		logging.CommonLog().Error("[Config] Error reading file: %s\n", err)
		return cfg, err
	}
	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		logging.CommonLog().Error("[Config] Error unmarshaling yaml file: %s\n", err)
		return cfg, err
	}
	cfg.SetDefaults()
	if err := cfg.Validate(); err != nil {
		logging.CommonLog().Error(err)
	}
	return cfg, err
}
