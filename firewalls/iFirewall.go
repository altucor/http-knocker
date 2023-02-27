package firewalls

import (
	"github.com/altucor/http-knocker/common"
	"github.com/altucor/http-knocker/devices"
	"github.com/altucor/http-knocker/firewallCommon/firewallField"
)

type IFirewall interface {
	GetDevice() devices.IDevice
	GetEndpoint() common.EndpointCfg
	AddClient(ip_addr firewallField.Address) error
	CleanupExpiredClients() error
}
