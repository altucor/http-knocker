package firewallProtocol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

func (ctx *routerOsRestRule) fromProtocol(data map[string]interface{}) (firewallCommon.FirewallRule, error) {
	// modify map keys and values to conform common vision
	rule := firewallCommon.FirewallRule{}
	rule.Id.TryInitFromString(data[".id"].(string))
	rule.Action.TryInitFromString(data["action"].(string))
	rule.Chain.TryInitFromString(data["chain"].(string))
	rule.Disabled.TryInitFromString(data["disabled"].(string))
	rule.Protocol.TryInitFromString(data["protocol"].(string))
	rule.SrcAddress.TryInitFromString(data["src-address"].(string))
	rule.DstPort.TryInitFromString(data["dst-port"].(string))
	rule.Comment.TryInitFromString(data["comment"].(string))
	rule.Detail.TryInitFromString(data["detail"].(string))
	rule.ErrorCmd.TryInitFromString(data["error"].(string))
	rule.Message.TryInitFromString(data["message"].(string))

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
	return "router-os-rest"
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
		method = http.MethodPut
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
		responseCmd = response.Get{}
		restProto := routerOsRestRule{}
		for _, element := range jsonArr {
			rule, err := restProto.fromProtocol(element)
			if err != nil {
				logging.CommonLog().Error("[FirewallRuleList] Error parsing firewall rule: %s\n", err)
				continue
			}
			responseCmd.(response.Get).AppendRule(rule)
		}
	case device.DeviceCommandRemove:
		responseCmd = response.Remove{}
	}
	return responseCmd, nil
}
