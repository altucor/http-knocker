package deviceCommon

import (
	"httpKnocker/firewallCommon"
	"net/http"
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
	Rest(http.Response)
}
