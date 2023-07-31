package firewallProtocol

import (
	"fmt"
	"strings"

	"github.com/altucor/http-knocker/device"
	"github.com/altucor/http-knocker/device/command"
	"github.com/altucor/http-knocker/device/response"
	"github.com/altucor/http-knocker/firewallCommon"
	"github.com/altucor/http-knocker/logging"
)

type routerOsTerminalRule struct {
}

func (ctx *routerOsTerminalRule) toProtocol(rule firewallCommon.FirewallRule) string {
	var result string = ""
	result += "place-before=" + rule.PlaceBefore.GetString()
	result += " action=" + rule.Action.GetString()
	result += " chain=" + rule.Chain.GetString()
	result += " disabled=" + rule.Detail.GetString()
	result += " protocol=" + rule.Protocol.GetString()
	result += " src-address=" + rule.SrcAddress.GetString()
	result += " dst-port=" + rule.DstPort.GetString()
	result += " comment=" + rule.Comment.GetString()
	return result
}

func (ctx *routerOsTerminalRule) fromProtocol(data string) (firewallCommon.FirewallRule, error) {
	// modify map keys and values to conform common vision
	rule := firewallCommon.FirewallRule{}

	// if err := initRuleField(&rule.Id, data, ".id"); err != nil {
	// 	return rule, err
	// }
	// if err := initRuleField(&rule.Action, data, "action"); err != nil {
	// 	return rule, err
	// }
	// if err := initRuleField(&rule.Chain, data, "chain"); err != nil {
	// 	return rule, err
	// }
	// if err := initRuleField(&rule.Disabled, data, "disabled"); err != nil {
	// 	return rule, err
	// }
	// if err := initRuleField(&rule.Protocol, data, "protocol"); err != nil {
	// 	return rule, err
	// }
	// if err := initRuleField(&rule.SrcAddress, data, "src-address"); err != nil {
	// 	return rule, err
	// }
	// if err := initRuleField(&rule.DstPort, data, "dst-port"); err != nil {
	// 	return rule, err
	// }
	// if err := initRuleField(&rule.Comment, data, "comment"); err != nil {
	// 	return rule, err
	// }
	// if err := initRuleField(&rule.Detail, data, "detail"); err != nil {
	// 	return rule, err
	// }
	// if err := initRuleField(&rule.ErrorCmd, data, "error"); err != nil {
	// 	return rule, err
	// }
	// if err := initRuleField(&rule.Message, data, "message"); err != nil {
	// 	return rule, err
	// }

	return rule, nil
}

type ProtocolRouterOsTerminal struct {
}

func (ctx ProtocolRouterOsTerminal) GetType() string {
	return "router-os-terminal"
}

func (ctx ProtocolRouterOsTerminal) To(cmd device.IDeviceCommand) (string, error) {
	var cmdData = ""
	switch cmd.GetType() {
	case device.DeviceCommandAdd:
		frw := routerOsTerminalRule{}
		cmdData = "/ip firewall add " + frw.toProtocol(cmd.(command.Add).GetRule())
	case device.DeviceCommandGet:
		cmdData = "/ip firewall filter print without-paging"
	case device.DeviceCommandRemove:
		cmdData = fmt.Sprintf("/ip firewall filter remove numbers=%d", cmd.(command.Remove).GetId())
	}
	return cmdData, nil
}

func (ctx ProtocolRouterOsTerminal) From(data string, cmdType device.DeviceCommandType) (device.IDeviceResponse, error) {
	var responseCmd device.IDeviceResponse

	switch cmdType {
	case device.DeviceCommandAdd:
		responseCmd = &response.Add{}
	case device.DeviceCommandGet:
		parts := strings.Split(data, "\n\r")
		getResponse := &response.Get{}
		ruleProto := routerOsTerminalRule{}
		for _, element := range parts {
			if strings.HasPrefix(element, "Flags: X - disabled, I - invalid, D - dynamic") {
				continue
			}
			if strings.HasPrefix(element, "# ") {
				continue
			}
			rule, err := ruleProto.fromProtocol(element)
			if err != nil {
				logging.CommonLog().Error("[ProtocolRouterOsRest] Error parsing firewall rule:", err)
				continue
			}
			getResponse.AppendRule(rule)
		}
		// That assign at end because can't figure out how to cast interface to pointer structure
		responseCmd = getResponse
	case device.DeviceCommandRemove:
		responseCmd = &response.Remove{}
	}
	return responseCmd, nil
}
