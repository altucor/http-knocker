package command

import (
	"github.com/altucor/http-knocker/device"
)

type Get struct {
	cmdType device.DeviceCommandType
}

func GetNew() Get {
	cmd := Get{
		cmdType: device.DeviceCommandGet,
	}
	return cmd
}

func (ctx Get) ToMap() map[string]interface{} {
	cmd := make(map[string]interface{})
	cmd["type"] = string(ctx.cmdType)
	return cmd
}

func (ctx Get) GetType() device.DeviceCommandType {
	return ctx.cmdType
}
