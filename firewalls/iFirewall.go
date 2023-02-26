package firewalls

import (
	"httpKnocker/common"
	"httpKnocker/devices"
	"httpKnocker/firewallCommon/firewallField"
)

type IFirewall interface {
	GetDevice() devices.IDevice
	GetEndpoint() common.EndpointCfg
	AddClient(ip_addr firewallField.Address) error
	CleanupExpiredClients() error
}
