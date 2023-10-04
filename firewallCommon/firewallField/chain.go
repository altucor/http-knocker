package firewallField

import (
	"errors"
	"strings"
)

type ChainType uint8

const (
	CHAIN_INVALID ChainType = 0
	CHAIN_INPUT   ChainType = 1
	CHAIN_FORWARD ChainType = 2
	CHAIN_OUTPUT  ChainType = 3
)

var (
	chainMap = map[ChainType]string{
		CHAIN_INVALID: "<INVALID>",
		CHAIN_INPUT:   "input",
		CHAIN_FORWARD: "forward",
		CHAIN_OUTPUT:  "output",
	}
)

type Chain struct {
	value ChainType
}

func (ctx *Chain) TryInitFromString(param string) error {
	if len(param) > 0 {
		param = strings.ToLower(param)
		for key, value := range chainMap {
			if value == param {
				ctx.value = key
				return nil
			}
		}
	}
	ctx.value = CHAIN_INVALID
	return errors.New("cannot init from string")
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
