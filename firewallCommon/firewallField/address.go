package firewallField

import (
	"net/netip"

	"github.com/altucor/http-knocker/logging"
)

type Address struct {
	value netip.Addr
}

func (ctx *Address) TryInitFromString(param string) error {
	addr, err := netip.ParseAddr(param)
	if err != nil {
		logging.CommonLog().Error("Cannot init from string, %s\n", err)
		return err
	}
	ctx.value = addr
	return nil
}

func AddressTypeFromString(idString string) (Address, error) {
	id := Address{}
	return id, id.TryInitFromString(idString)
}

func (ctx *Address) SetValue(value netip.Addr) {
	ctx.value = value
}

func (ctx Address) GetValue() netip.Addr {
	return ctx.value
}

func (ctx Address) GetString() string {
	return ctx.value.String()
}
