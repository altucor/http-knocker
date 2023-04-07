package response

import (
	"github.com/altucor/http-knocker/device"
)

type Remove struct {
	cmdType device.DeviceCommandType
	err     error
}

func (ctx Remove) GetType() device.DeviceCommandType {
	return ctx.cmdType
}

func (ctx *Remove) SetError(err error) {
	ctx.err = err
}

func (ctx Remove) GetError() error {
	return ctx.err
}
