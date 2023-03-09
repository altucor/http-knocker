package deviceCommon

import (
	"github.com/altucor/http-knocker/firewallCommon"
)

type DeviceCommandType string

const (
	DeviceCommandGet    DeviceCommandType = "get"
	DeviceCommandAdd    DeviceCommandType = "add"
	DeviceCommandMove   DeviceCommandType = "move"
	DeviceCommandRemove DeviceCommandType = "remove"
)

type IDeviceCommand interface {
	GetType() DeviceCommandType
	Rest() (string, string, string, error) // Return: method, url, body
	IpTables() (string, error)
}

type IDeviceResponse interface {
	GetType() DeviceCommandType
	GetRules() firewallCommon.FirewallRuleList
}
