package firewallCommon

import (
	"github.com/altucor/http-knocker/firewallCommon/firewallField"
)

const RULE_ID_INVALID = 0xFFFFFFFFFFFFFFFF

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

func (ctx *FirewallRule) ToMap() map[string]string {
	ruleMap := make(map[string]string)
	if ctx.Id.GetValue() != RULE_ID_INVALID {
		ruleMap["id"] = ctx.Id.GetString()
	}
	ruleMap["action"] = ctx.Action.GetString()
	ruleMap["chain"] = ctx.Chain.GetString()
	ruleMap["disabled"] = ctx.Disabled.GetString()
	ruleMap["protocol"] = ctx.Protocol.GetString()
	ruleMap["src-address"] = ctx.SrcAddress.GetString()
	ruleMap["dst-port"] = ctx.DstPort.GetString()
	ruleMap["comment"] = ctx.Comment.GetString()
	ruleMap["place-before"] = ctx.PlaceBefore.GetString()

	return ruleMap
}

func (ctx *FirewallRule) FromMap(m map[string]string) {
	if _, ok := m["id"]; ok {
		ctx.Id.TryInitFromString(m["id"])
	}
	if _, ok := m["action"]; ok {
		ctx.Action.TryInitFromString(m["action"])
	}
	if _, ok := m["chain"]; ok {
		ctx.Chain.TryInitFromString(m["chain"])
	}
	if _, ok := m["disabled"]; ok {
		ctx.Disabled.TryInitFromString(m["disabled"])
	}
	if _, ok := m["protocol"]; ok {
		ctx.Protocol.TryInitFromString(m["protocol"])
	}
	if _, ok := m["src-address"]; ok {
		ctx.SrcAddress.TryInitFromString(m["src-address"])
	}
	if _, ok := m["dst-port"]; ok {
		ctx.DstPort.TryInitFromString(m["dst-port"])
	}
	if _, ok := m["comment"]; ok {
		ctx.Comment.TryInitFromString(m["comment"])
	}
	if _, ok := m["place-before"]; ok {
		ctx.PlaceBefore.TryInitFromString(m["place-before"])
	}
}
