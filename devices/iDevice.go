package devices

import (
	"errors"

	"github.com/altucor/http-knocker/deviceCommon"
	"github.com/altucor/http-knocker/logging"
	"golang.org/x/exp/slices"
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

type DeviceConnectionDesc struct {
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Endpoint   string `yaml:"endpoint"`
	Tls        bool   `yaml:"tls"`
	Host       string `yaml:"host"`
	Port       uint16 `yaml:"port"`
	KnownHosts string `yaml:"knownHosts"`
}

type IDevice interface {
	Start() error
	Stop() error
	GetSupportedProtocols() []DeviceProtocol
	GetType() DeviceType
	RunCommandWithReply(cmd deviceCommon.IDeviceCommand, proto DeviceProtocol) (deviceCommon.IDeviceResponse, error)
	// RunCommand(firewallCommon.IFirewallCommand) error
}
