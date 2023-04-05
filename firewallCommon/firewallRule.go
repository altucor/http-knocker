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
