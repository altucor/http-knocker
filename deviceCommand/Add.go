package deviceCommand

import (
	"net/http"
	"net/netip"

	"github.com/altucor/http-knocker/deviceCommon"
	"github.com/altucor/http-knocker/firewallCommon"
	"github.com/altucor/http-knocker/firewallCommon/firewallField"
	"github.com/altucor/http-knocker/logging"
)

type Add struct {
	cmdType      deviceCommon.DeviceCommandType
	firewallRule firewallCommon.FirewallRule
}

func AddNew(clientAddr netip.Addr, port uint16, protocol firewallField.ProtocolType, comment string, placeBefore uint64) Add {
	frwRule := firewallCommon.FirewallRuleNew()
	frwRule.Id.SetValue(firewallCommon.RULE_ID_INVALID)
	frwRule.Action.SetValue(firewallField.ACCEPT)
	frwRule.Chain.SetValue(firewallField.INPUT)
	frwRule.Disabled.SetValue(false)
	frwRule.Protocol.SetValue(protocol)
	frwRule.SrcAddress.SetValue(clientAddr)
	frwRule.DstPort.SetValue(port)
	frwRule.Comment.SetValue(comment)
	frwRule.PlaceBefore.SetValue(placeBefore)

	cmd := Add{
		cmdType:      deviceCommon.DeviceCommandAdd,
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

func (ctx Add) GetType() deviceCommon.DeviceCommandType {
	return ctx.cmdType
}

func (ctx Add) GetRule() firewallCommon.FirewallRule {
	return ctx.firewallRule
}

func (ctx Add) Rest() (string, string, string, error) {
	method := http.MethodPut
	url := "/ip/firewall/filter"
	body, err := ctx.firewallRule.ToRest()
	if err != nil {
		logging.CommonLog().Error("[deviceCommandAdd] Error converting firewall rule to REST: %s\n", err)
		return "", "", "", err
	}

	return method, url, body, nil
}

func (ctx Add) IpTables() (string, error) {
	return ctx.firewallRule.ToIpTables()
}
