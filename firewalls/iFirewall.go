package firewalls

import (
	"http-knocker/common"
	"http-knocker/devices"
	"http-knocker/firewallCommon/firewallField"
)

type IFirewall interface {
	GetDevice() devices.IDevice
	GetEndpoint() common.EndpointCfg
	AddClient(ip_addr firewallField.Address) error
	CleanupExpiredClients() error
}
