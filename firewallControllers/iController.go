package firewallControllers

import (
	"net/http"

	"github.com/altucor/http-knocker/devices"
	"github.com/altucor/http-knocker/endpoint"
	"gopkg.in/yaml.v3"
)

type IController interface {
	SetDevice(dev devices.IDevice)
	SetEndpoint(endpoint *endpoint.Endpoint)
	Start() error
	Stop() error
	HttpCallbackAddClient(w http.ResponseWriter, r *http.Request)
	GetHttpCallback() (string, func(w http.ResponseWriter, r *http.Request))
	CleanupExpiredClients() error
}

type Config struct {
	Type     string `yaml:"type"`
	Device   string `yaml:"device"`
	Endpoint string `yaml:"endpoint"`
}

type InterfaceWrapper struct {
	Controller IController
	Config     Config
}

func (ctx *InterfaceWrapper) UnmarshalYAML(value *yaml.Node) error {
	if err := value.Decode(&ctx.Config); err != nil {
		return err
	}

	var err error = nil
	switch ctx.Config.Type {
	case "basic":
		ctx.Controller, err = ControllerBasicNewFromYaml(value)
	}
	if err != nil {
		return err
	}
	return nil
}
