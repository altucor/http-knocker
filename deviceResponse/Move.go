package deviceResponse

import (
	"net/http"

	"github.com/altucor/http-knocker/deviceCommon"
	"github.com/altucor/http-knocker/firewallCommon"
)

type Move struct {
	cmdType deviceCommon.DeviceCommandType
	rules   firewallCommon.FirewallRuleList
}

func MoveFromRouterOsRest(response http.Response) (Move, error) {
	frwResponse := Move{
		cmdType: deviceCommon.DeviceCommandMove,
	}
	rules, err := firewallCommon.FirewallRuleListNewFromRest(response)
	if err != nil {
		return Move{}, err
	}
	frwResponse.rules = rules
	return frwResponse, nil
}

func MoveFromIpTables(response string) (Move, error) {
	frwResponse := Move{
		cmdType: deviceCommon.DeviceCommandMove,
	}
	rules, err := firewallCommon.FirewallRuleListNewFromIpTables(response)
	if err != nil {
		return Move{}, err
	}
	frwResponse.rules = rules
	return frwResponse, nil
}

func (ctx Move) GetType() deviceCommon.DeviceCommandType {
	return ctx.cmdType
}

func (ctx Move) GetRules() firewallCommon.FirewallRuleList {
	return ctx.rules
}
