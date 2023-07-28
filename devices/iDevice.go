package devices

import (
	"github.com/altucor/http-knocker/device"
	"github.com/altucor/http-knocker/firewallProtocol"
)

type IDevice interface {
	SetProtocol(protocol firewallProtocol.IFirewallProtocol)
	Start() error
	Stop() error
	RunCommandWithReply(cmd device.IDeviceCommand) (device.IDeviceResponse, error)
}
