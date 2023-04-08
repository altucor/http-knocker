package devices

import (
	"errors"

	"github.com/altucor/http-knocker/device"
	"github.com/altucor/http-knocker/firewallProtocol"
	"github.com/altucor/http-knocker/logging"
	"gopkg.in/yaml.v3"
)

type IDevice interface {
	Start() error
	Stop() error
	RunCommandWithReply(cmd device.IDeviceCommand) (device.IDeviceResponse, error)
}

type Config struct {
	Type     string `yaml:"type"`
	Protocol string `yaml:"protocol"`
}

type InterfaceWrapper struct {
	Device IDevice
	Config Config
}

func (ctx *InterfaceWrapper) UnmarshalYAML(value *yaml.Node) error {
	if err := value.Decode(&ctx.Config); err != nil {
		return err
	}

	protocolStorage := firewallProtocol.GetProtocolStorage()
	protocol := protocolStorage.GetProtocolByName(ctx.Config.Protocol)

	var err error = nil
	switch ctx.Config.Type {
	case "rest":
		ctx.Device, err = DeviceRestNewFromYaml(value, protocol.(IFirewallRestProtocol))
	case "ssh":
		ctx.Device, err = DeviceSshNewFromYaml(value, protocol.(IFirewallSshProtocol))
	case "puller":
		ctx.Device, err = DevicePullerNewFromYaml(value, nil)
	case "router-os":
		ctx.Device, err = DeviceRouterOsNewFromYaml(value, nil)
	default:
		logging.CommonLog().Error("[iDevice] invalid type of device")
		return errors.New("invalid type of device")
	}
	if err != nil {
		return err
	}
	return nil
}
