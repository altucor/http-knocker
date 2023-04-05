package response

import (
	"github.com/altucor/http-knocker/device"
)

type Add struct {
	cmdType device.DeviceCommandType
}

func (ctx Add) GetType() device.DeviceCommandType {
	return ctx.cmdType
}
