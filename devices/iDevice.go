package devices

import (
	"errors"

	"github.com/altucor/http-knocker/device"
	"github.com/altucor/http-knocker/firewallProtocol"
	"github.com/altucor/http-knocker/logging"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

type DeviceType string

const (
	DeviceTypeSsh      DeviceType = "ssh"
	DeviceTypeRouterOs DeviceType = "routeros"
	DeviceTypeRest     DeviceType = "rest"
	DeviceTypePuller   DeviceType = "puller"
)

var (
	deviceTypeArr = []DeviceType{
		DeviceTypeSsh,
		DeviceTypeRouterOs,
		DeviceTypeRest,
		DeviceTypePuller,
	}
)

func (ctx *DeviceType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tempStr string = ""
	err := unmarshal(&tempStr)
	if err != nil {
		return err
	}
	tempDevType := DeviceType(tempStr)
	if !slices.Contains(deviceTypeArr, tempDevType) {
		logging.CommonLog().Error("Cannot init from string")
		return errors.New("cannot init from string")
	}
	*ctx = tempDevType
	return nil
}

type IDevice interface {
	Start() error
	Stop() error
	GetType() DeviceType
	RunCommandWithReply(cmd device.IDeviceCommand) (device.IDeviceResponse, error)
}

type InterfaceWrapper struct {
	Device   IDevice
	Type     string `yaml:"type"`
	Protocol string `yaml:"protocol"`
}

func (ctx *InterfaceWrapper) UnmarshalYAML(value *yaml.Node) error {
	var intermediate InterfaceWrapper
	if err := value.Decode(&intermediate); err != nil {
		return err
	}
	ctx = &intermediate

	var err error = nil
	var protocol firewallProtocol.IFirewallProtocol
	switch intermediate.Protocol {
	case "rest-router-os":
		protocol = firewallProtocol.ProtocolRouterOsRest{}
	case "ssh-iptables":
		protocol = firewallProtocol.ProtocolIpTables{}
	case "puller":
		protocol = nil
	default:
		logging.CommonLog().Error("[iDevice] invalid type of protocol")
		return errors.New("invalid type of protocol")
	}

	switch intermediate.Type {
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
