package deviceCommand

import (
	"httpKnocker/deviceCommon"
	"httpKnocker/firewallCommon"
	"httpKnocker/firewallCommon/firewallField"
	"httpKnocker/logging"
	"net/http"
	"net/netip"
)

type Add struct {
	cmdType         deviceCommon.DeviceCommandType
	firewallRule    firewallCommon.FirewallRule
	durationSeconds uint64
}

func AddNew(clientAddr netip.Addr, port uint16, protocol firewallField.ProtocolType, durationSeconds uint64, comment string, placeBefore uint64) Add {
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
		cmdType:         deviceCommon.DeviceCommandAdd,
		firewallRule:    frwRule,
		durationSeconds: durationSeconds,
	}
	return cmd
}

func (ctx Add) GetType() deviceCommon.DeviceCommandType {
	return ctx.cmdType
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
