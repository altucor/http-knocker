package firewallField

import (
	"fmt"
	"strconv"

	"github.com/altucor/http-knocker/logging"
)

type Port struct {
	value uint16
}

func (ctx *Port) TryInitFromString(param string) error {
	port, err := strconv.ParseUint(param, 10, 16)
	if err != nil {
		logging.CommonLog().Error("Cannot init from string, %s\n", err)
		return err
	}
	ctx.value = uint16(port)
	return nil
}

func (ctx *Port) TryInitFromRest(param string) error {
	return ctx.TryInitFromString(param)
}

func (ctx *Port) TryInitFromIpTables(param string) error {
	return ctx.TryInitFromString(param)
}

func PortTypeFromString(idString string) (Port, error) {
	id := Port{}
	return id, id.TryInitFromString(idString)
}

func PortTypeFromValue(value uint16) (Port, error) {
	return Port{value: value}, nil
}

func (ctx *Port) SetValue(value uint16) {
	ctx.value = value
}

func (ctx Port) GetValue() uint16 {
	return ctx.value
}

func (ctx Port) GetString() string {
	return fmt.Sprintf("%d", ctx.value)
}

func (ctx Port) MarshalRest() string {
	return fmt.Sprintf("%d", ctx.value)
}

func (ctx Port) MarshalIpTables() string {
	return ctx.GetString()
}
