package response

import (
	"github.com/altucor/http-knocker/device"
)

type Remove struct {
	cmdType device.DeviceCommandType
}

func (ctx Remove) GetType() device.DeviceCommandType {
	return ctx.cmdType
}
