package firewallField

import (
	"errors"
	"strconv"
)

type Bool struct {
	value bool
}

func BoolNew(val bool) Bool {
	return Bool{value: val}
}

func (ctx *Bool) TryInitFromString(param string) error {
	value, err := strconv.ParseBool(param)
	if err != nil {
		ctx.value = false
		return errors.New("cannot init from string")
	}
	ctx.value = value
	return nil
}

func (ctx *Bool) SetValue(value bool) {
	ctx.value = value
}

func (ctx Bool) GetValue() bool {
	return ctx.value
}

func (ctx Bool) GetString() string {
	if ctx.value {
		return "true"
	}
	return "false"
}
