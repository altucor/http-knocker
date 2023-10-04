package firewallCommon

import (
	"github.com/altucor/http-knocker/firewallCommon/firewallField"
)

// const RULE_ID_INVALID = 0xFFFFFFFFFFFFFFFF

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
	if ctx.Id.GetValue() != firewallField.RULE_ID_INVALID {
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

func (ctx *FirewallRule) FromMap(m map[string]string) error {
	if _, ok := m["id"]; ok {
		if err := ctx.Id.TryInitFromString(m["id"]); err != nil {
			return err
		}
	}
	if _, ok := m["action"]; ok {
		if err := ctx.Action.TryInitFromString(m["action"]); err != nil {
			return err
		}
	}
	if _, ok := m["chain"]; ok {
		if err := ctx.Chain.TryInitFromString(m["chain"]); err != nil {
			return err
		}
	}
	if _, ok := m["disabled"]; ok {
		if err := ctx.Disabled.TryInitFromString(m["disabled"]); err != nil {
			return err
		}
	}
	if _, ok := m["protocol"]; ok {
		if err := ctx.Protocol.TryInitFromString(m["protocol"]); err != nil {
			return err
		}
	}
	if _, ok := m["src-address"]; ok {
		if err := ctx.SrcAddress.TryInitFromString(m["src-address"]); err != nil {
			return err
		}
	}
	if _, ok := m["dst-port"]; ok {
		if err := ctx.DstPort.TryInitFromString(m["dst-port"]); err != nil {
			return err
		}
	}
	if _, ok := m["comment"]; ok {
		if err := ctx.Comment.TryInitFromString(m["comment"]); err != nil {
			return err
		}
	}
	if _, ok := m["place-before"]; ok {
		if err := ctx.PlaceBefore.TryInitFromString(m["place-before"]); err != nil {
			return err
		}
	}

	return nil
}

func FirewallRuleFromMap(m map[string]string) FirewallRule {
	f := FirewallRule{}
	f.FromMap(m)
	return f
}
