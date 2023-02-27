package firewallCommon

import (
	"encoding/json"
	"fmt"
	"http-knocker/firewallCommon/firewallField"
	"http-knocker/logging"
	"regexp"
)

type FirewallRule struct {
	Id          firewallField.Number
	Action      firewallField.Action
	Chain       firewallField.Chain
	Disabled    firewallField.Bool
	Protocol    firewallField.Protocol
	SrcAddress  firewallField.Address
	DstPort     firewallField.Port
	Comment     firewallField.Text
	Detail      firewallField.Text
	ErrorCmd    firewallField.Number
	Message     firewallField.Text
	PlaceBefore firewallField.Number
}

func FirewallRuleNew() FirewallRule {
	frwRule := FirewallRule{
		Id:          firewallField.Number{},
		Action:      firewallField.Action{},
		Chain:       firewallField.Chain{},
		Disabled:    firewallField.Bool{},
		Protocol:    firewallField.Protocol{},
		SrcAddress:  firewallField.Address{},
		DstPort:     firewallField.Port{},
		Comment:     firewallField.Text{},
		Detail:      firewallField.Text{},
		ErrorCmd:    firewallField.Number{},
		Message:     firewallField.Text{},
		PlaceBefore: firewallField.Number{},
	}
	return frwRule
}

func FirewallRuleNewFromRestMap(data map[string]string) (FirewallRule, error) {
	frwRule := FirewallRuleNew()
	for key, element := range data {
		switch key {
		case ".id":
			frwRule.Id.TryInitFromRest(element)
		case "action":
			frwRule.Action.TryInitFromRest(element)
		case "chain":
			frwRule.Chain.TryInitFromRest(element)
		case "disabled":
			frwRule.Disabled.TryInitFromRest(element)
		case "protocol":
			frwRule.Protocol.TryInitFromRest(element)
		case "src-address":
			frwRule.SrcAddress.TryInitFromRest(element)
		case "dst-port":
			frwRule.DstPort.TryInitFromRest(element)
		case "comment":
			frwRule.Comment.TryInitFromRest(element)
		case "detail":
			frwRule.Detail.TryInitFromRest(element)
		case "error":
			frwRule.ErrorCmd.TryInitFromRest(element)
		case "message":
			frwRule.Message.TryInitFromRest(element)
		}
	}

	return frwRule, nil
}

func iptablesParseParamViaRegex(rule string, frwParam IFirewallField, re string) {
	res := regexp.MustCompile(re).FindStringSubmatch(rule)
	if len(res) < 2 {
		return
	}
	if res[1] != "" {
		frwParam.TryInitFromIpTables(res[1])
	}
}

func FirewallRuleNewFromIpTables(rule string) (FirewallRule, error) {
	// Regex strings from: https://www.npmjs.com/package/@wkronmiller/iptables-parser?activeTab=explore
	frwRule := FirewallRuleNew()
	iptablesParseParamViaRegex(rule, &frwRule.Chain, `-A\s([^\s]+)\s`)
	iptablesParseParamViaRegex(rule, &frwRule.Action, `\s-j\s([^\s]+)`)
	iptablesParseParamViaRegex(rule, &frwRule.Protocol, `\s-p\s([A-Za-z]+)`)
	iptablesParseParamViaRegex(rule, &frwRule.SrcAddress, `\s-s\s([^\s]+)`)
	iptablesParseParamViaRegex(rule, &frwRule.DstPort, `\s--dport\s([^\s]+)`)
	iptablesParseParamViaRegex(rule, &frwRule.Comment, `-m\s+comment\s+--comment\s+("[^"]*"|'[^']*'|[^'"\s]+)`)

	return frwRule, nil
}

func (ctx *FirewallRule) ToRest() (string, error) {
	jsonMap := make(map[string]string)
	if ctx.Id.GetValue() != RULE_ID_INVALID {
		jsonMap[".id"] = ctx.Id.MarshalRest()
	}
	jsonMap["action"] = ctx.Action.MarshalRest()
	jsonMap["chain"] = ctx.Chain.MarshalRest()
	jsonMap["disabled"] = ctx.Disabled.MarshalRest()
	jsonMap["protocol"] = ctx.Protocol.MarshalRest()
	jsonMap["src-address"] = ctx.SrcAddress.MarshalRest()
	jsonMap["dst-port"] = ctx.DstPort.MarshalRest()
	jsonMap["comment"] = ctx.Comment.MarshalRest()
	jsonMap["place-before"] = ctx.PlaceBefore.MarshalRest()

	jsonBytes, err := json.Marshal(jsonMap)
	if err != nil {
		logging.CommonLog().Error("[FirewallRule] Error marshaling to json: %s\n", err)
		return "", err
	}

	return string(jsonBytes), nil
}

func (ctx *FirewallRule) ToIpTables() (string, error) {
	// iptables -I INPUT 2 -p tcp -s 10.1.1.2 --dport 22 -j ACCEPT -m comment --comment "My comments here"
	var result string = "iptables -t filter "
	result += "-I " + ctx.Chain.MarshalIpTables() + " "
	if ctx.PlaceBefore.GetValue() != RULE_ID_INVALID {
		result += fmt.Sprintf("%d", ctx.PlaceBefore.GetValue()) + " "
	}
	result += "-p " + ctx.Protocol.MarshalIpTables() + " "
	result += "-s " + ctx.SrcAddress.MarshalIpTables() + " "
	result += "--dport " + fmt.Sprintf("%d", ctx.DstPort.GetValue()) + " "
	result += "-j " + ctx.Action.MarshalIpTables() + " "
	result += "-m comment --comment \"" + ctx.Comment.MarshalIpTables() + "\""

	return result, nil
}
