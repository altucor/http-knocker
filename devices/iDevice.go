package devices

import (
	"github.com/altucor/http-knocker/deviceCommon"
)

type DeviceType string

const (
	DeviceTypeSsh      DeviceType = "ssh"
	DeviceTypeRouterOs DeviceType = "routeros"
	DeviceTypeRest     DeviceType = "rest"
	DeviceTypePuller   DeviceType = "puller"
)

type IDevice interface {
	Start() error
	Stop() error
	GetSupportedProtocols() []DeviceProtocol
	GetType() DeviceType
	RunCommandWithReply(cmd deviceCommon.IDeviceCommand, proto DeviceProtocol) (deviceCommon.IDeviceResponse, error)
	// RunCommand(firewallCommon.IFirewallCommand) error
}
