package firewallProtocol

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/altucor/http-knocker/device"
	"github.com/altucor/http-knocker/device/command"
	"github.com/altucor/http-knocker/device/response"
	"github.com/altucor/http-knocker/firewallCommon"
	"github.com/altucor/http-knocker/logging"
)

type IpTablesRule struct {
}

func (ctx *IpTablesRule) toProtocol(rule firewallCommon.FirewallRule) string {
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

func (ctx *IpTablesRule) fromProtocol(data string) (firewallCommon.FirewallRule, error) {
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

type ProtocolIpTables struct {
}

func (ctx ProtocolIpTables) GetType() string {
	return "ssh-iptables"
}

func (ctx ProtocolIpTables) To(cmd device.IDeviceCommand) (string, error) {
	var cmdData = ""
	switch cmd.GetType() {
	case device.DeviceCommandAdd:
		frw := IpTablesRule{}
		cmdData = "iptables -t filter " + frw.toProtocol(cmd.(command.Add).GetRule())
	case device.DeviceCommandGet:
		cmdData = "iptables -S INPUT"
	case device.DeviceCommandRemove:
		// TODO: Re-Check removing indexes.
		// Because for PullerDevice numeration of rules starts from 1, NOT FROM 0
		cmdData = fmt.Sprintf("iptables --delete INPUT %d", cmd.(command.Remove).GetId()+1)
	}
	return cmdData, nil
}

func (ctx ProtocolIpTables) From(data string, cmdType device.DeviceCommandType) (device.IDeviceResponse, error) {
	var responseCmd device.IDeviceResponse

	if len(data) == 0 {
		return responseCmd, nil
	}

	switch cmdType {
	case device.DeviceCommandAdd:
		responseCmd = &response.Add{}
	case device.DeviceCommandGet:
		getResponse := &response.Get{}
		rules := strings.Split(data, "\r\n")
		ruleProto := IpTablesRule{}
		for index, element := range rules {
			if len(element) == 0 {
				continue
			}
			rule, err := ruleProto.fromProtocol(element)
			if err != nil {
				logging.CommonLog().Error("[FirewallRuleList] Error parsing firewall rule:", err)
				return responseCmd, err
			}
			rule.Id.SetValue(uint64(index))
			getResponse.AppendRule(rule)
		}
		responseCmd = getResponse
	case device.DeviceCommandRemove:
		responseCmd = &response.Remove{}
	}
	return responseCmd, nil
}
