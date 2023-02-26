package deviceResponse

import (
	"httpKnocker/deviceCommon"
	"httpKnocker/firewallCommon"
	"net/http"
)

type Add struct {
	cmdType deviceCommon.DeviceCommandType
	rules   firewallCommon.FirewallRuleList
}

func AddFromRouterOsRest(response http.Response) (Add, error) {
	frwResponse := Add{
		cmdType: deviceCommon.DeviceCommandAdd,
	}
	rules, err := firewallCommon.FirewallRuleListNewFromRest(response)
	if err != nil {
		return Add{}, err
	}
	frwResponse.rules = rules
	return frwResponse, nil
}

func AddFromIpTables(response string) (Add, error) {
	frwResponse := Add{
		cmdType: deviceCommon.DeviceCommandAdd,
	}
	rules, err := firewallCommon.FirewallRuleListNewFromIpTables(response)
	if err != nil {
		return Add{}, err
	}
	frwResponse.rules = rules
	return frwResponse, nil
}

func (ctx Add) GetType() deviceCommon.DeviceCommandType {
	return ctx.cmdType
}

func (ctx Add) GetRules() firewallCommon.FirewallRuleList {
	return ctx.rules
}

func (ctx Add) Rest(http.Response) {

}
