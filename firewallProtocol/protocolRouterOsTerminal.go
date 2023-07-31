package firewallProtocol

import (
	"github.com/altucor/http-knocker/device"
	"github.com/altucor/http-knocker/device/response"
	"github.com/altucor/http-knocker/firewallCommon"
)

type routerOsTerminalRule struct {
}

func (ctx *routerOsTerminalRule) toProtocol(rule firewallCommon.FirewallRule) string {
	return ""
}

func (ctx *routerOsTerminalRule) fromProtocol(data string) (firewallCommon.FirewallRule, error) {
	return firewallCommon.FirewallRule{}, nil
}

type ProtocolRouterOsTerminal struct {
}

func (ctx ProtocolRouterOsTerminal) GetType() string {
	return "router-os-terminal"
}

func (ctx ProtocolRouterOsTerminal) To(cmd device.IDeviceCommand) (string, error) {
	return "", nil
}

func (ctx ProtocolRouterOsTerminal) From(data string, cmdType device.DeviceCommandType) (device.IDeviceResponse, error) {
	return &response.Add{}, nil
}
