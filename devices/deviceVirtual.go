package devices

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/altucor/http-knocker/device"
	"github.com/altucor/http-knocker/device/command"
	"github.com/altucor/http-knocker/device/response"
	"github.com/altucor/http-knocker/firewallCommon"
)

type virtualFirewallCmd struct {
	id  string
	cmd device.IDeviceCommand
}

func (ctx virtualFirewallCmd) toMap() map[string]interface{} {
	vfcmd := make(map[string]interface{})
	vfcmd["id"] = ctx.id
	vfcmd["command"] = ctx.cmd.ToMap()
	return vfcmd
}

type VirtualFirewall struct {
	mu    sync.Mutex
	cmds  []virtualFirewallCmd
	rules []firewallCommon.FirewallRule
}

/*
Important:
Instead of storing rules in to virtual firewall.
Save actual commands in virtual firewall state.
Because communication with remote client should be with commands, not rules
Remote client should know how to interpret commands in rules
and how to execute commands for custom firewall
*/

func generateCmdId() string {
	h := sha1.New()
	h.Write([]byte(time.Now().String()))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (ctx *VirtualFirewall) Add(cmd command.Add) (*response.Add, error) {
	// Should add new commands to pending list until they will be accepted by remote firewall
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.cmds = append(ctx.cmds, virtualFirewallCmd{
		id:  generateCmdId(),
		cmd: cmd,
	})
	return &response.Add{}, nil
}

func (ctx *VirtualFirewall) Get(cmd command.Get) (*response.Get, error) {
	// Should return with list of accepted virtual firewall rules
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	return response.GetFromRuleList(ctx.rules)
}

func (ctx *VirtualFirewall) Remove(cmd command.Remove) (*response.Remove, error) {
	// Should mark rules from accepted list as pending for removal, but not remove them
	// Only really remove them when remote firewall will approve this
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.cmds = append(ctx.cmds, virtualFirewallCmd{
		id:  generateCmdId(),
		cmd: cmd,
	})
	return &response.Remove{}, nil
}

func (ctx *VirtualFirewall) getLastPendingCommands(count uint64) (string, error) {
	// Here we respond only with pending changes for remote firewall
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	var cmds []map[string]interface{}
	for _, item := range ctx.cmds {
		if count != 0 && uint64(len(cmds)) >= count {
			break
		}
		cmds = append(cmds, item.toMap())
	}
	jsonBytes, err := json.Marshal(cmds)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func (ctx *VirtualFirewall) processAcceptedCommands(acceptedCommands []string) error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	for _, acceptedCommand := range acceptedCommands {
		for i, item := range ctx.cmds {
			if item.id == acceptedCommand {
				// Removing this command from pending list
				// No need to emulate execution of command
				// Result of execution should be received from remote device
				ctx.cmds = append(ctx.cmds[:i], ctx.cmds[i+1:]...)
				break
			}
		}
	}
	return nil
}

func (ctx *VirtualFirewall) pushRuleSet(rules []firewallCommon.FirewallRule) error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.rules = rules
	return nil
}
