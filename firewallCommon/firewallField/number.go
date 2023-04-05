package firewallField

import (
	"errors"
	"fmt"
	"strconv"
)

type Number struct {
	value uint64
}

func (ctx *Number) TryInitFromString(param string) error {
	value, err := strconv.ParseUint(param, 16, 64)
	if err != nil {
		ctx.value = 0xFFFFFFFFFFFFFFFF
		return errors.New("cannot init from string")
	}
	ctx.value = value
	return nil
}

func NumberTypeFromString(idString string) (Number, error) {
	id := Number{}
	return id, id.TryInitFromString(idString)
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
