package firewallControllers

import (
	"errors"
)

type ControllerValueType uint32

const (
	FIREWALL_INVALID ControllerValueType = 0xFFFFFFFF
	FIREWALL_BASIC   ControllerValueType = 0
)

var (
	firewallsMap = map[ControllerValueType]string{
		FIREWALL_INVALID: "<INVALID>",
		FIREWALL_BASIC:   "firewallBasic",
	}
)

type ControllerType struct {
	value ControllerValueType
}

func (ctx *ControllerType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var firewallStr string = ""
	err := unmarshal(&firewallStr)
	if err != nil {
		return err
	}
	firewall, err := ControllerTypeFromString(firewallStr)
	if err != nil {
		return err
	}
	ctx.SetValue(firewall.GetValue())
	return nil
}

func ControllerTypeFromString(firewallString string) (ControllerType, error) {
	for key, value := range firewallsMap {
		if value == firewallString {
			return ControllerType{value: key}, nil
		}
	}

	return ControllerType{value: FIREWALL_INVALID}, errors.New("invalid firewall text name")
}

func ControllerTypeFromValue(value ControllerValueType) (ControllerType, error) {
	_, ok := firewallsMap[value]
	if ok {
		return ControllerType{value: value}, nil
	}
	return ControllerType{value: FIREWALL_INVALID}, errors.New("invalid furewall type value")
}

func (ctx *ControllerType) SetValue(value ControllerValueType) {
	ctx.value = value
}

func (ctx ControllerType) GetValue() ControllerValueType {
	return ctx.value
}

func (ctx ControllerType) GetString() string {
	return firewallsMap[ctx.value]
}
