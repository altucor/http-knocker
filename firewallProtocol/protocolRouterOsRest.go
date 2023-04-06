package firewallProtocol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/altucor/http-knocker/device"
	"github.com/altucor/http-knocker/device/command"
	"github.com/altucor/http-knocker/device/response"
	"github.com/altucor/http-knocker/firewallCommon"
	"github.com/altucor/http-knocker/logging"
)

type routerOsRestRule struct {
}

func (ctx *routerOsRestRule) toProtocol(rule firewallCommon.FirewallRule) map[string]interface{} {
	// modify map keys and values to conform device protocol
	mapObj := make(map[string]interface{})
	if rule.Id.GetValue() != firewallCommon.RULE_ID_INVALID {
		mapObj[".id"] = fmt.Sprintf("*%X", rule.Id.GetValue())
	}
	mapObj["action"] = rule.Action.GetString()
	mapObj["chain"] = rule.Chain.GetString()
	mapObj["disabled"] = rule.Disabled.GetString()
	mapObj["protocol"] = rule.Protocol.GetString()
	mapObj["src-address"] = rule.SrcAddress.GetString()
	mapObj["dst-port"] = rule.DstPort.GetString()
	mapObj["comment"] = rule.Comment.GetString()
	mapObj["place-before"] = rule.PlaceBefore.GetString()
	return mapObj
}

func initRuleField(field firewallCommon.IFirewallField, data map[string]interface{}, name string) error {
	if _, ok := data[name]; ok {
		if name == ".id" {
			if err := field.TryInitFromString(strings.ReplaceAll(data[name].(string), "*", "")); err != nil {
				return err
			}
		} else {
			if err := field.TryInitFromString(data[name].(string)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (ctx *routerOsRestRule) fromProtocol(data map[string]interface{}) (firewallCommon.FirewallRule, error) {
	// modify map keys and values to conform common vision
	rule := firewallCommon.FirewallRule{}

	if err := initRuleField(&rule.Id, data, ".id"); err != nil {
		return rule, err
	}
	if err := initRuleField(&rule.Action, data, "action"); err != nil {
		return rule, err
	}
	if err := initRuleField(&rule.Chain, data, "chain"); err != nil {
		return rule, err
	}
	if err := initRuleField(&rule.Disabled, data, "disabled"); err != nil {
		return rule, err
	}
	if err := initRuleField(&rule.Protocol, data, "protocol"); err != nil {
		return rule, err
	}
	if err := initRuleField(&rule.SrcAddress, data, "src-address"); err != nil {
		return rule, err
	}
	if err := initRuleField(&rule.DstPort, data, "dst-port"); err != nil {
		return rule, err
	}
	if err := initRuleField(&rule.Comment, data, "comment"); err != nil {
		return rule, err
	}
	if err := initRuleField(&rule.Detail, data, "detail"); err != nil {
		return rule, err
	}
	if err := initRuleField(&rule.ErrorCmd, data, "error"); err != nil {
		return rule, err
	}
	if err := initRuleField(&rule.Message, data, "message"); err != nil {
		return rule, err
	}

	return rule, nil
}

// func FirewallRuleListNewFromRest(response http.Response) ([]FirewallRule, error) {
// 	frwList := make([]FirewallRule, 1)

// 	body, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		logging.CommonLog().Error("[FirewallRuleList] Error reading response body: %s\n", err)
// 		return frwList, err
// 	}

// 	if len(body) == 0 {
// 		return frwList, nil
// 	}

// 	if strings.HasPrefix(string(body), "{") {
// 		// Single item
// 		// TODO: Fix decoding of non string value from json
// 		// {"detail":"unknown parameter .id","error":400,"message\":\"
// 		// var testDecode map[string]interface{}
// 		// testType := testDecode["detail"].(type)
// 		// err = json.Unmarshal(body, &testDecode)
// 		// logging.CommonLog().Debugf("test decode: %s\n", testDecode)
// 		var jsonMap map[string]string
// 		err = json.Unmarshal(body, &jsonMap)
// 		if err != nil {
// 			logging.CommonLog().Error("[FirewallRuleList] Error unmarshal json to map: %s\n", err)
// 			return frwList, err
// 		}
// 		// logging.DebugLogger.Printf("[FirewallRuleList] jsonMap: %s\n", jsonMap)
// 		rule, err := FirewallRuleNewFromRestMap(jsonMap)
// 		if err != nil {
// 			logging.CommonLog().Error("[FirewallRuleList] Error parsing firewall rule: %s\n", err)
// 			return frwList, err
// 		}
// 		frwList = append(frwList, rule)
// 	} else if strings.HasPrefix(string(body), "[") {
// 		// Multiple items
// 		var jsonArr []map[string]string
// 		err = json.Unmarshal([]byte(body), &jsonArr)
// 		if err != nil {
// 			logging.CommonLog().Error("[FirewallRuleList] Error unmarshal json to array: %s\n", err)
// 			return frwList, err
// 		}
// 		// logging.DebugLogger.Printf("[FirewallRuleList] jsonArr: %s\n", jsonArr)
// 		for _, element := range jsonArr {
// 			rule, err := FirewallRuleNewFromRestMap(element)
// 			if err != nil {
// 				logging.CommonLog().Error("[FirewallRuleList] Error parsing firewall rule: %s\n", err)
// 				continue
// 			}
// 			frwList = append(frwList, rule)
// 		}
// 	}

// 	return frwList, nil
// }

type ProtocolRouterOsRest struct {
}

func (ctx ProtocolRouterOsRest) GetType() FirewallProtocolName {
	return "rest-router-os"
}

func (ctx ProtocolRouterOsRest) To(cmd device.IDeviceCommand, baseUrl string) (*http.Request, error) {
	var method string = ""
	var url string = "/ip/firewall/filter"
	var body string = ""

	switch cmd.GetType() {
	case device.DeviceCommandAdd:
		method = http.MethodPut
		restProto := routerOsRestRule{}
		jsonBytes, err := json.Marshal(restProto.toProtocol(cmd.(command.Add).GetRule()))
		if err != nil {
			logging.CommonLog().Error("[ProtocolRouterOsRest] Error marshaling to json: %s\n", err)
			return nil, err
		}
		body = string(jsonBytes)
	case device.DeviceCommandGet:
		method = http.MethodGet
	case device.DeviceCommandRemove:
		method = http.MethodDelete
		url += fmt.Sprintf("/*%X", cmd.(command.Remove).GetId())
	}

	req, err := http.NewRequest(method, baseUrl+url, bytes.NewReader([]byte(body)))
	if err != nil {
		logging.CommonLog().Error("could not create request: %s\n", err)
	}
	req.Header.Set("content-type", "application/json")
	return req, nil
}

func (ctx ProtocolRouterOsRest) From(
	httpResponse *http.Response,
	cmdType device.DeviceCommandType,
) (device.IDeviceResponse, error) {
	var responseCmd device.IDeviceResponse

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		logging.CommonLog().Error("[FirewallRuleList] Error reading response body: %s\n", err)
		return responseCmd, err
	}
	if len(body) == 0 {
		return responseCmd, nil
	}
	var jsonArr []map[string]interface{}
	err = json.Unmarshal([]byte(body), &jsonArr)
	if err != nil {
		logging.CommonLog().Error("[FirewallRuleList] Error unmarshal json to array: %s\n", err)
		return responseCmd, err
	}

	switch cmdType {
	case device.DeviceCommandAdd:
		responseCmd = response.Add{}
	case device.DeviceCommandGet:
		getResponse := response.Get{}
		restProto := routerOsRestRule{}
		for _, element := range jsonArr {
			rule, err := restProto.fromProtocol(element)
			if err != nil {
				logging.CommonLog().Error("[FirewallRuleList] Error parsing firewall rule: %s\n", err)
				continue
			}
			getResponse.AppendRule(rule)
		}
		// That assign at end because can't figure out how to cast interface to pointer structure
		responseCmd = getResponse
	case device.DeviceCommandRemove:
		responseCmd = &response.Remove{}
	}
	return responseCmd, nil
}
