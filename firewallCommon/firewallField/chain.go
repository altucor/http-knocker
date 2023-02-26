package firewallField

import (
	"errors"
	"strings"
)

type ChainType uint8

const (
	INVALID ChainType = 0xFF
	INPUT   ChainType = 0
	FORWARD ChainType = 1
	OUTPUT  ChainType = 2
)

var (
	chainMap = map[ChainType]string{
		INVALID: "<INVALID>",
		INPUT:   "input",
		FORWARD: "forward",
		OUTPUT:  "output",
	}
)

type Chain struct {
	value ChainType
}

func (ctx *Chain) TryInitFromString(param string) error {
	param = strings.ToLower(param)
	for key, value := range chainMap {
		if value == param {
			ctx.value = key
			return nil
		}
	}
	ctx.value = INVALID
	return errors.New("Cannot init from string")
}

func (ctx *Chain) TryInitFromRest(param string) error {
	return ctx.TryInitFromString(param)
}

func (ctx *Chain) TryInitFromIpTables(param string) error {
	return ctx.TryInitFromString(param)
}

func ChainTypeFromString(chainString string) (Chain, error) {
	chainString = strings.ToLower(chainString)
	for key, value := range chainMap {
		if value == chainString {
			return Chain{value: key}, nil
		}
	}

	return Chain{value: INVALID}, errors.New("Invalid chain text name")
}

func ChainTypeFromValue(value ChainType) (Chain, error) {
	_, ok := chainMap[value]
	if ok {
		return Chain{value: value}, nil
	}
	return Chain{value: INVALID}, errors.New("Invalid protocol type value")
}

func (ctx *Chain) SetValue(value ChainType) {
	ctx.value = value
}

func (ctx Chain) GetValue() ChainType {
	return ctx.value
}

func (ctx Chain) GetString() string {
	return chainMap[ctx.value]
}

func (ctx Chain) MarshalRest() string {
	return chainMap[ctx.value]
}

func (ctx Chain) MarshalIpTables() string {
	return strings.ToUpper(chainMap[ctx.value])
}
