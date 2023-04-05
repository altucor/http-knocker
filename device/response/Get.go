package response

import (
	"github.com/altucor/http-knocker/device"
	"github.com/altucor/http-knocker/firewallCommon"
)

type Get struct {
	cmdType device.DeviceCommandType
	rules   []firewallCommon.FirewallRule
}

// func GetFromRouterOsRest(response http.Response) (Get, error) {
// 	frwResponse := Get{
// 		cmdType: deviceCommon.DeviceCommandGet,
// 	}

// 	rules, err := firewallCommon.FirewallRuleListNewFromRest(response)
// 	if err != nil {
// 		return Get{}, err
// 	}
// 	frwResponse.rules = rules
// 	return frwResponse, nil
// }

// func GetFromIpTables(response string) (Get, error) {
// 	frwResponse := Get{
// 		cmdType: deviceCommon.DeviceCommandGet,
// 	}

// 	rules, err := firewallCommon.FirewallRuleListNewFromIpTables(response)
// 	if err != nil {
// 		return Get{}, err
// 	}
// 	frwResponse.rules = rules
// 	return frwResponse, nil
// }

func GetFromRuleList(rules []firewallCommon.FirewallRule) (Get, error) {
	frwResponse := Get{
		cmdType: device.DeviceCommandGet,
		rules:   rules,
	}

	return frwResponse, nil
}

func (ctx Get) GetType() device.DeviceCommandType {
	return ctx.cmdType
}

func (ctx Get) GetRules() []firewallCommon.FirewallRule {
	return ctx.rules
}

func (ctx Get) AppendRule(rule firewallCommon.FirewallRule) {
	ctx.rules = append(ctx.rules, rule)
}
