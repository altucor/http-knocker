package firewallField

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const RULE_ID_INVALID = 0xFFFFFFFFFFFFFFFF

type Number struct {
	value uint64
}

func NumberNew(val uint64) Number {
	return Number{value: val}
}

func (ctx *Number) TryInitFromString(param string) error {
	param = strings.ReplaceAll(param, "0x", "")
	value, err := strconv.ParseUint(param, 16, 64)
	if err != nil {
		ctx.value = RULE_ID_INVALID
		return errors.New("cannot init from string")
	}
	ctx.value = value
	return nil
}

func (ctx *Number) SetValue(value uint64) {
	ctx.value = value
}

func (ctx Number) GetValue() uint64 {
	return ctx.value
}

func (ctx Number) GetString() string {
	return fmt.Sprintf("0x%02X", ctx.value)
}
