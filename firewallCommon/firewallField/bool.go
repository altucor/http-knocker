package firewallField

import (
	"errors"
	"strconv"
)

type Bool struct {
	value bool
}

func (ctx *Bool) TryInitFromString(param string) error {
	value, err := strconv.ParseBool(param)
	if err != nil {
		ctx.value = false
		return errors.New("Cannot init from string")
	}
	ctx.value = value
	return nil
}

func BoolTypeFromString(idString string) (Bool, error) {
	id := Bool{}
	return id, id.TryInitFromString(idString)
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
