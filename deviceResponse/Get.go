package deviceResponse

import (
	"net/http"

	"github.com/altucor/http-knocker/deviceCommon"
	"github.com/altucor/http-knocker/firewallCommon"
)

type Get struct {
	cmdType deviceCommon.DeviceCommandType
	rules   firewallCommon.FirewallRuleList
}

func GetFromRouterOsRest(response http.Response) (Get, error) {
	frwResponse := Get{
		cmdType: deviceCommon.DeviceCommandGet,
	}

	rules, err := firewallCommon.FirewallRuleListNewFromRest(response)
	if err != nil {
		return Get{}, err
	}
	frwResponse.rules = rules
	return frwResponse, nil
}

func GetFromIpTables(response string) (Get, error) {
	frwResponse := Get{
		cmdType: deviceCommon.DeviceCommandGet,
	}

	rules, err := firewallCommon.FirewallRuleListNewFromIpTables(response)
	if err != nil {
		return Get{}, err
	}
	frwResponse.rules = rules
	return frwResponse, nil
}

func (ctx Get) GetType() deviceCommon.DeviceCommandType {
	return ctx.cmdType
}

func (ctx Get) GetRules() firewallCommon.FirewallRuleList {
	return ctx.rules
}
