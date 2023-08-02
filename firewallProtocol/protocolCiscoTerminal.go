package firewallProtocol

import (
	"github.com/altucor/http-knocker/device"
	"github.com/altucor/http-knocker/firewallCommon"
)

type ciscoTerminalRule struct {
}

func (ctx *ciscoTerminalRule) toProtocol(rule firewallCommon.FirewallRule) string {
	return ""
}

func (ctx *ciscoTerminalRule) fromProtocol(data string) (firewallCommon.FirewallRule, error) {
	return firewallCommon.FirewallRule{}, nil
}

type ProtocolCiscoTerminal struct {
}

func (ctx ProtocolCiscoTerminal) GetType() string {
	return "cisco-terminal"
}

func (ctx ProtocolCiscoTerminal) To(cmd device.IDeviceCommand) (string, error) {
	return "", nil
}

func (ctx ProtocolCiscoTerminal) From(data string, cmdType device.DeviceCommandType) (device.IDeviceResponse, error) {
	return nil, nil
}
