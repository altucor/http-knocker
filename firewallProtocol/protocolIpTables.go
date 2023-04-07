package firewallProtocol

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/altucor/http-knocker/device"
	"github.com/altucor/http-knocker/device/command"
	"github.com/altucor/http-knocker/device/response"
	"github.com/altucor/http-knocker/firewallCommon"
)

type IpTablesRule struct {
}

func (ctx *IpTablesRule) ToProtocol(rule firewallCommon.FirewallRule) string {
	// modify map keys and values to conform device protocol
	// iptables -I INPUT 2 -p tcp -s 10.1.1.2 --dport 22 -j ACCEPT -m comment --comment "My comments here"
	var result string = ""
	result += "-I " + strings.ToUpper(rule.Chain.GetString()) + " "
	if rule.PlaceBefore.GetValue() != firewallCommon.RULE_ID_INVALID {
		result += fmt.Sprintf("%d", rule.PlaceBefore.GetValue()) + " "
	}
	result += "-p " + rule.Protocol.GetString() + " "
	result += "-s " + rule.SrcAddress.GetString() + " "
	result += "--dport " + fmt.Sprintf("%d", rule.DstPort.GetValue()) + " "
	result += "-j " + strings.ToUpper(rule.Action.GetString()) + " "
	result += "-m comment --comment \"" + rule.Comment.GetString() + "\""
	return result
}

func iptablesParseParamViaRegex(rule string, frwParam firewallCommon.IFirewallField, re string) {
	res := regexp.MustCompile(re).FindStringSubmatch(rule)
	if len(res) < 2 {
		return
	}
	if res[1] != "" {
		frwParam.TryInitFromString(res[1])
	}
}

func (ctx *IpTablesRule) FromProtocol(data string) (firewallCommon.FirewallRule, error) {
	// modify map keys and values to conform common vision
	rule := firewallCommon.FirewallRule{}
	iptablesParseParamViaRegex(data, &rule.Chain, `-A\s([^\s]+)\s`)
	iptablesParseParamViaRegex(data, &rule.Action, `\s-j\s([^\s]+)`)
	iptablesParseParamViaRegex(data, &rule.Protocol, `\s-p\s([A-Za-z]+)`)
	iptablesParseParamViaRegex(data, &rule.SrcAddress, `\s-s\s([^\s]+)`)
	iptablesParseParamViaRegex(data, &rule.DstPort, `\s--dport\s([^\s]+)`)
	iptablesParseParamViaRegex(data, &rule.Comment, `-m\s+comment\s+--comment\s+("[^"]*"|'[^']*'|[^'"\s]+)`)
	return rule, nil
}

// func FirewallRuleListNewFromIpTables(response string) ([]FirewallRule, error) {
// 	frwList := make([]FirewallRule, 1)

// 	if len(response) == 0 {
// 		return frwList, nil
// 	}

// 	rules := strings.Split(response, "\r\n")
// 	for index, element := range rules {
// 		rule, err := FirewallRuleNewFromIpTables(element)
// 		if err != nil {
// 			logging.CommonLog().Error("[FirewallRuleList] Error parsing firewall rule: %s\n", err)
// 			return frwList, err
// 		}
// 		rule.Id.SetValue(uint64(index))
// 		frwList = append(frwList, rule)
// 	}

// 	return frwList, nil
// }

type ProtocolIpTables struct {
}

func (ctx ProtocolIpTables) GetType() FirewallProtocolName {
	return "ssh-iptables"
}

func (ctx ProtocolIpTables) To(cmd device.IDeviceCommand) (string, error) {
	var cmdData = ""
	switch cmd.GetType() {
	case device.DeviceCommandAdd:
		frw := IpTablesRule{}
		cmdData = "iptables -t filter " + frw.ToProtocol(cmd.(command.Add).GetRule())
	case device.DeviceCommandGet:
		cmdData = "iptables -S INPUT"
	case device.DeviceCommandRemove:
		cmdData = fmt.Sprintf("iptables --delete INPUT %d", cmd.(command.Remove).GetId())
	}
	return cmdData, nil
}

func (ctx ProtocolIpTables) From(data string, cmdType device.DeviceCommandType) (device.IDeviceResponse, error) {
	var responseCmd device.IDeviceResponse
	switch cmdType {
	case device.DeviceCommandAdd:
		responseCmd = response.Add{}
	case device.DeviceCommandGet:
		responseCmd = response.Get{}
	case device.DeviceCommandRemove:
		responseCmd = response.Remove{}
	}
	return responseCmd, nil
}