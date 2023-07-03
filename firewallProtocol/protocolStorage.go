package firewallProtocol

import "sync"

type ProtocolStorage struct {
	protocols map[string]IFirewallProtocol
}

func (ctx *ProtocolStorage) Init() {
	ctx.protocols = make(map[string]IFirewallProtocol)
	ctx.protocols["rest-router-os"] = ProtocolRouterOsRest{}
	ctx.protocols["ssh-iptables"] = ProtocolIpTables{}
}

func (ctx *ProtocolStorage) GetProtocolByName(name string) IFirewallProtocol {
	return ctx.protocols[name]
}

var lock = &sync.Mutex{}
var protocolStorage *ProtocolStorage

func GetProtocolStorage() *ProtocolStorage {
	lock.Lock()
	defer lock.Unlock()
	if protocolStorage == nil {
		protocolStorage = &ProtocolStorage{}
		protocolStorage.Init()
	}

	return protocolStorage
}
