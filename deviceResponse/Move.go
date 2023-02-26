package deviceResponse

import (
	"httpKnocker/deviceCommon"
	"httpKnocker/firewallCommon"
	"net/http"
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

func (ctx Move) Rest(http.Response) {

}
