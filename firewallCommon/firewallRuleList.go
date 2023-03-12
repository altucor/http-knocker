package firewallCommon

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/altucor/http-knocker/logging"
)

func FirewallRuleListNewFromRest(response http.Response) ([]FirewallRule, error) {
	frwList := make([]FirewallRule, 1)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logging.CommonLog().Error("[FirewallRuleList] Error reading response body: %s\n", err)
		return frwList, err
	}

	if len(body) == 0 {
		return frwList, nil
	}

	if strings.HasPrefix(string(body), "{") {
		// Single item
		// TODO: Fix decoding of non string value from json
		// {"detail":"unknown parameter .id","error":400,"message\":\"
		// var testDecode map[string]interface{}
		// testType := testDecode["detail"].(type)
		// err = json.Unmarshal(body, &testDecode)
		// logging.CommonLog().Debugf("test decode: %s\n", testDecode)
		var jsonMap map[string]string
		err = json.Unmarshal(body, &jsonMap)
		if err != nil {
			logging.CommonLog().Error("[FirewallRuleList] Error unmarshal json to map: %s\n", err)
			return frwList, err
		}
		// logging.DebugLogger.Printf("[FirewallRuleList] jsonMap: %s\n", jsonMap)
		rule, err := FirewallRuleNewFromRestMap(jsonMap)
		if err != nil {
			logging.CommonLog().Error("[FirewallRuleList] Error parsing firewall rule: %s\n", err)
			return frwList, err
		}
		frwList = append(frwList, rule)
	} else if strings.HasPrefix(string(body), "[") {
		// Multiple items
		var jsonArr []map[string]string
		err = json.Unmarshal([]byte(body), &jsonArr)
		if err != nil {
			logging.CommonLog().Error("[FirewallRuleList] Error unmarshal json to array: %s\n", err)
			return frwList, err
		}
		// logging.DebugLogger.Printf("[FirewallRuleList] jsonArr: %s\n", jsonArr)
		for _, element := range jsonArr {
			rule, err := FirewallRuleNewFromRestMap(element)
			if err != nil {
				logging.CommonLog().Error("[FirewallRuleList] Error parsing firewall rule: %s\n", err)
				continue
			}
			frwList = append(frwList, rule)
		}
	}

	return frwList, nil
}

func FirewallRuleListNewFromIpTables(response string) ([]FirewallRule, error) {
	frwList := make([]FirewallRule, 1)

	if len(response) == 0 {
		return frwList, nil
	}

	rules := strings.Split(response, "\r\n")
	for index, element := range rules {
		rule, err := FirewallRuleNewFromIpTables(element)
		if err != nil {
			logging.CommonLog().Error("[FirewallRuleList] Error parsing firewall rule: %s\n", err)
			return frwList, err
		}
		rule.Id.SetValue(uint64(index))
		frwList = append(frwList, rule)
	}

	return frwList, nil
}
