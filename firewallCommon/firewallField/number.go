package firewallField

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Number struct {
	value uint64
}

func (ctx *Number) TryInitFromString(param string) error {
	value, err := strconv.ParseUint(param, 16, 64)
	if err != nil {
		ctx.value = 0xFFFFFFFFFFFFFFFF
		return errors.New("Cannot init from string")
	}
	ctx.value = value
	return nil
}

func (ctx *Number) TryInitFromRest(param string) error {
	return ctx.TryInitFromString(strings.ReplaceAll(param, "*", ""))
}

func NumberTypeFromString(idString string) (Number, error) {
	id := Number{}
	return id, id.TryInitFromString(idString)
}

func NumberTypeFromValue(value uint64) (Number, error) {
	return Number{value: value}, nil
}

func (ctx *Number) SetValue(value uint64) {
	ctx.value = value
}

func (ctx Number) GetValue() uint64 {
	return ctx.value
}

func (ctx Number) GetString() string {
	return fmt.Sprintf("%d", ctx.value)
}

func (ctx Number) MarshalRest() string {
	return fmt.Sprintf("*%X", ctx.value)
}

func (ctx Number) MarshalIpTables() string {
	return ctx.GetString()
}
