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

type InterfaceWrapper struct {
	Controller IController
	Type       string `yaml:"type"`
	Device     string `yaml:"device"`
	Endpoint   string `yaml:"endpoint"`
}

func (ctx *InterfaceWrapper) UnmarshalYAML(value *yaml.Node) error {
	var intermediate InterfaceWrapper
	if err := value.Decode(&intermediate); err != nil {
		return err
	}
	ctx = &intermediate

	var err error = nil
	switch intermediate.Type {
	case "basic":
		ctx.Controller, err = ControllerBasicNewFromYaml(value)
	}
	if err != nil {
		return err
	}
	return nil
}
