package response

import (
	"github.com/altucor/http-knocker/device"
)

type Add struct {
	cmdType device.DeviceCommandType
	err     error
}

func (ctx Add) GetType() device.DeviceCommandType {
	return ctx.cmdType
}

func (ctx *Add) SetError(err error) {
	ctx.err = err
}

func (ctx Add) GetError() error {
	return ctx.err
}
