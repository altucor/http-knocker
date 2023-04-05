package command

import (
	"github.com/altucor/http-knocker/device"
)

type Remove struct {
	cmdType device.DeviceCommandType
	ruleId  uint64
}

func RemoveNew(id uint64) Remove {
	frw := Remove{
		cmdType: device.DeviceCommandRemove,
		ruleId:  id,
	}
	return frw
}

func (ctx Remove) ToMap() map[string]interface{} {
	cmd := make(map[string]interface{})
	cmd["type"] = string(ctx.cmdType)
	cmd["id"] = ctx.ruleId
	return cmd
}

func (ctx Remove) GetType() device.DeviceCommandType {
	return ctx.cmdType
}

func (ctx Remove) GetId() uint64 {
	return ctx.ruleId
}
