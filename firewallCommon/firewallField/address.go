package firewallField

import (
	"errors"
	"http-knocker/logging"
	"net/netip"
	"strings"
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

func (ctx *Address) TryInitFromRest(param string) error {
	return ctx.TryInitFromString(param)
}

func (ctx *Address) TryInitFromIpTables(param string) error {
	parts := strings.Split(param, "/")
	if len(parts) != 2 {
		logging.CommonLog().Error("Cannot detect mask delimiter in CIDR string")
		return errors.New("Cannot detect mask delimiter in CIDR string")
	}
	return ctx.TryInitFromString(parts[0])
}

func AddressTypeFromString(idString string) (Address, error) {
	id := Address{}
	return id, id.TryInitFromString(idString)
}

func AddressTypeFromValue(value netip.Addr) (Address, error) {
	return Address{value: value}, nil
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

func (ctx Address) MarshalRest() string {
	return ctx.GetString()
}

func (ctx Address) MarshalIpTables() string {
	return ctx.GetString()
}
