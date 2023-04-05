package device

type DeviceCommandType string

const (
	DeviceCommandGet    DeviceCommandType = "get"
	DeviceCommandAdd    DeviceCommandType = "add"
	DeviceCommandRemove DeviceCommandType = "remove"
)

type IDeviceCommand interface {
	ToMap() map[string]interface{}
	GetType() DeviceCommandType
}

type IDeviceResponse interface {
	GetType() DeviceCommandType
}
