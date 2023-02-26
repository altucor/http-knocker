package deviceResponse

import (
	"httpKnocker/deviceCommon"
	"httpKnocker/firewallCommon"
	"net/http"
)

type Remove struct {
	cmdType deviceCommon.DeviceCommandType
	rules   firewallCommon.FirewallRuleList
}

func RemoveFromRouterOsRest(response http.Response) (Remove, error) {
	frwResponse := Remove{
		cmdType: deviceCommon.DeviceCommandRemove,
	}
	rules, err := firewallCommon.FirewallRuleListNewFromRest(response)
	if err != nil {
		return Remove{}, err
	}
	frwResponse.rules = rules
	return frwResponse, nil
}

func RemoveFromIpTables(response string) (Remove, error) {
	frwResponse := Remove{
		cmdType: deviceCommon.DeviceCommandRemove,
	}
	rules, err := firewallCommon.FirewallRuleListNewFromIpTables(response)
	if err != nil {
		return Remove{}, err
	}
	frwResponse.rules = rules
	return frwResponse, nil
}

func (ctx Remove) GetType() deviceCommon.DeviceCommandType {
	return ctx.cmdType
}

func (ctx Remove) GetRules() firewallCommon.FirewallRuleList {
	return ctx.rules
}

func (ctx Remove) Rest(http.Response) {

}
