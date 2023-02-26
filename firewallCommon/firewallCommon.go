package firewallCommon

import (
	"errors"
)

const RULE_ID_INVALID = 0xFFFFFFFFFFFFFFFF

type FirewallValueType uint32

const (
	FIREWALL_INVALID FirewallValueType = 0xFFFFFFFF
	FIREWALL_BASIC   FirewallValueType = 0
	FIREWALL_PULL    FirewallValueType = 1
)

var (
	firewallsMap = map[FirewallValueType]string{
		FIREWALL_INVALID: "<INVALID>",
		FIREWALL_BASIC:   "firewallBasic",
		FIREWALL_PULL:    "firewallPull",
	}
)

type FirewallType struct {
	value FirewallValueType
}

func (ctx *FirewallType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var firewallStr string = ""
	err := unmarshal(&firewallStr)
	if err != nil {
		return err
	}
	firewall, err := FirewallTypeFromString(firewallStr)
	if err != nil {
		return err
	}
	ctx.SetValue(firewall.GetValue())
	return nil
}

func FirewallTypeFromString(firewallString string) (FirewallType, error) {
	for key, value := range firewallsMap {
		if value == firewallString {
			return FirewallType{value: key}, nil
		}
	}

	return FirewallType{value: FIREWALL_INVALID}, errors.New("Invalid firewall text name")
}

func FirewallTypeFromValue(value FirewallValueType) (FirewallType, error) {
	_, ok := firewallsMap[value]
	if ok {
		return FirewallType{value: value}, nil
	}
	return FirewallType{value: FIREWALL_INVALID}, errors.New("Invalid furewall type value")
}

func (ctx FirewallType) SetValue(value FirewallValueType) {
	ctx.value = value
}

func (ctx FirewallType) GetValue() FirewallValueType {
	return ctx.value
}

func (ctx FirewallType) GetString() string {
	return firewallsMap[ctx.value]
}
