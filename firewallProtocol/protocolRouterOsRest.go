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
	"github.com/altucor/http-knocker/firewallCommon/firewallField"
	"github.com/altucor/http-knocker/logging"
)

type routerOsRestRule struct {
}

func (ctx *routerOsRestRule) toProtocol(rule firewallCommon.FirewallRule) map[string]interface{} {
	// modify map keys and values to conform device protocol
	mapObj := make(map[string]interface{})
	if rule.Id.GetValue() != firewallField.RULE_ID_INVALID {
		mapObj[".id"] = fmt.Sprintf("*%X", rule.Id.GetValue())
	}
	mapObj["action"] = rule.Action.GetString()
	mapObj["chain"] = rule.Chain.GetString()
	mapObj["disabled"] = rule.Disabled.GetString()
	mapObj["protocol"] = rule.Protocol.GetString()
	mapObj["src-address"] = rule.SrcAddress.GetString()
	mapObj["dst-port"] = rule.DstPort.GetString()
	mapObj["comment"] = rule.Comment.GetString()
	mapObj["place-before"] = fmt.Sprintf("*%X", rule.PlaceBefore.GetValue())
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

type ProtocolRouterOsRest struct {
}

func (ctx ProtocolRouterOsRest) GetType() string {
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
			logging.CommonLog().Error("[ProtocolRouterOsRest] Error marshaling to json:", err)
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
		logging.CommonLog().Error("could not create request:", err)
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
		logging.CommonLog().Error("[ProtocolRouterOsRest] Error reading response body:", err)
		return responseCmd, err
	}
	if len(body) == 0 {
		return responseCmd, nil
	}
	if http.StatusBadRequest <= httpResponse.StatusCode {
		var jsonErrorData map[string]interface{}
		err = json.Unmarshal([]byte(body), &jsonErrorData)
		if err != nil {
			logging.CommonLog().Error("[ProtocolRouterOsRest] Error unmarshal json to dict:", err)
			return responseCmd, err
		}
		return responseCmd, fmt.Errorf(
			"got error from device, code: %s message: %s",
			jsonErrorData["error"],
			jsonErrorData["message"],
		)
	}

	switch cmdType {
	case device.DeviceCommandAdd:
		responseCmd = &response.Add{}
	case device.DeviceCommandGet:
		var jsonArr []map[string]interface{}
		err = json.Unmarshal([]byte(body), &jsonArr)
		if err != nil {
			logging.CommonLog().Error("[ProtocolRouterOsRest] Error unmarshal json to array:", err)
			return responseCmd, err
		}
		getResponse := &response.Get{}
		restProto := routerOsRestRule{}
		for _, element := range jsonArr {
			rule, err := restProto.fromProtocol(element)
			if err != nil {
				logging.CommonLog().Error("[ProtocolRouterOsRest] Error parsing firewall rule:", err)
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
