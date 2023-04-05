package firewallControllers

import (
	"github.com/altucor/http-knocker/devices"
	"github.com/altucor/http-knocker/endpoint"
	"github.com/altucor/http-knocker/firewallCommon/firewallField"
	"gopkg.in/yaml.v3"
)

type IController interface {
	GetDeviceName() string
	SetDevice(dev devices.IDevice)
	GetEndpointName() string
	SetEndpoint(endpoint *endpoint.Endpoint)
	Start() error
	Stop() error
	GetDevice() devices.IDevice
	GetEndpoint() endpoint.Endpoint
	AddClient(ip_addr firewallField.Address) error
	CleanupExpiredClients() error
}

type InterfaceWrapper struct {
	Controller IController
}

func (ctx *InterfaceWrapper) UnmarshalYAML(value *yaml.Node) error {
	var intermediate struct {
		Type string `yaml:"type"`
	}
	if err := value.Decode(&intermediate); err != nil {
		return err
	}

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
