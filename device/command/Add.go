package command

import (
	"net/netip"

	"github.com/altucor/http-knocker/device"
	"github.com/altucor/http-knocker/firewallCommon"
	"github.com/altucor/http-knocker/firewallCommon/firewallField"
)

type Add struct {
	cmdType      device.DeviceCommandType
	firewallRule firewallCommon.FirewallRule
}

func AddNew(
	clientAddr netip.Addr,
	port uint16,
	protocol firewallField.ProtocolType,
	comment string,
	placeBefore uint64,
) Add {
	frwRule := firewallCommon.FirewallRule{}
	frwRule.Id.SetValue(firewallField.RULE_ID_INVALID)
	frwRule.Action.SetValue(firewallField.ACTION_ACCEPT)
	frwRule.Chain.SetValue(firewallField.CHAIN_INPUT)
	frwRule.Disabled.SetValue(false)
	frwRule.Protocol.SetValue(protocol)
	frwRule.SrcAddress.SetValue(clientAddr)
	frwRule.DstPort.SetValue(port)
	frwRule.Comment.SetValue(comment)
	frwRule.PlaceBefore.SetValue(placeBefore)

	cmd := Add{
		cmdType:      device.DeviceCommandAdd,
		firewallRule: frwRule,
	}
	return cmd
}

func (ctx Add) ToMap() map[string]interface{} {
	cmd := make(map[string]interface{})
	cmd["type"] = string(ctx.cmdType)
	cmd["rule"] = ctx.firewallRule.ToMap()
	return cmd
}

func (ctx Add) GetType() device.DeviceCommandType {
	return ctx.cmdType
}

func (ctx Add) GetRule() firewallCommon.FirewallRule {
	return ctx.firewallRule
}
