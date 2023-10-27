package firewallField

import (
	"fmt"
	"strconv"

	"github.com/altucor/http-knocker/logging"
)

type Port struct {
	value uint16
}

func PortNew(val uint16) Port {
	return Port{value: val}
}

func (ctx *Port) TryInitFromString(param string) error {
	port, err := strconv.ParseUint(param, 10, 16)
	if err != nil {
		logging.CommonLog().Errorf("cannot init from string, %s\n", err)
		return err
	}
	ctx.value = uint16(port)
	return nil
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
