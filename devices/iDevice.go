package devices

import (
	"github.com/altucor/http-knocker/device"
)

type IDevice interface {
	Start() error
	Stop() error
	RunCommandWithReply(cmd device.IDeviceCommand) (device.IDeviceResponse, error)
}
