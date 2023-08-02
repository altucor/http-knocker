package firewallProtocol

import (
	"sync"

	"github.com/altucor/http-knocker/logging"
)

type ProtocolStorage struct {
	protocols map[string]IFirewallProtocol
}

func (ctx *ProtocolStorage) Init() {
	ctx.protocols = make(map[string]IFirewallProtocol)
	ctx.protocols["rest-router-os"] = ProtocolRouterOsRest{}
	ctx.protocols["ssh-iptables"] = ProtocolIpTables{}
	ctx.protocols["puller"] = nil
	ctx.protocols["terminal-router-os"] = ProtocolRouterOsTerminal{}
	ctx.protocols["terminal-cisco"] = ProtocolCiscoTerminal{}
	ctx.protocols[""] = nil
}

func (ctx *ProtocolStorage) GetProtocolByName(name string) IFirewallProtocol {
	if _, ok := ctx.protocols[name]; !ok {
		logging.CommonLog().Fatalf("[ProtocolStorage] Cannot find protocol under name: \"%s\"", name)
	}
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
