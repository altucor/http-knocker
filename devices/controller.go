package devices

import (
	"errors"
	"fmt"
	"sync"

	"github.com/altucor/http-knocker/comment"
	"github.com/altucor/http-knocker/device"
	"github.com/altucor/http-knocker/device/command"
	"github.com/altucor/http-knocker/device/response"
	"github.com/altucor/http-knocker/firewallCommon"
	"github.com/altucor/http-knocker/firewallProtocol"
	"github.com/altucor/http-knocker/logging"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Type     string `yaml:"type"`
	Protocol string `yaml:"protocol"`
}

type DeviceConstructor interface {
	Construct(*yaml.Node, firewallProtocol.IFirewallProtocol) (IDevice, error)
}

type DeviceController struct {
	deviceMutex     sync.Mutex
	Device          IDevice
	Config          Config
	needCacheUpdate bool
	rulesCache      []firewallCommon.FirewallRule
}

func DeviceControllerNew(device IDevice) *DeviceController {
	ctx := DeviceController{
		Device: nil,
		Config: Config{
			Type:     "",
			Protocol: "",
		},
		needCacheUpdate: true,
	}

	return &ctx
}

func (ctx *DeviceController) UnmarshalYAML(value *yaml.Node) error {
	if err := value.Decode(&ctx.Config); err != nil {
		return err
	}

	protocolStorage := firewallProtocol.GetProtocolStorage()
	protocol := protocolStorage.GetProtocolByName(ctx.Config.Protocol)

	deviceStorage := GetDeviceStorage()
	device, err := deviceStorage.GetDevice(ctx.Config.Type, value)
	if err != nil {
		return err
	}
	device.SetProtocol(protocol)
	ctx.Device = device
	return nil
}

func (ctx *DeviceController) RemoveBatch(ids []uint64) error {
	for _, id := range ids {
		_, err := ctx.RunCommandWithReply(command.RemoveNew(id))
		if err != nil {
			return err
		}
	}
	return nil
}

func (ctx *DeviceController) GetRules(forceRefresh bool) ([]firewallCommon.FirewallRule, error) {
	// TODO: For now always update cache
	// In future update only after specific commands like ADD/REMOVE
	// and after specific interval
	ctx.needCacheUpdate = true
	if ctx.needCacheUpdate || forceRefresh {
		getResponse, err := ctx.RunCommandWithReply(command.GetNew())
		if err != nil {
			return nil, err
		}
		ctx.rulesCache = getResponse.(*response.Get).GetRules()
	}

	return ctx.rulesCache, nil
}

func (ctx *DeviceController) FindRuleIdByComment(comment string, forceRefresh bool) (uint64, error) {
	frwRules, err := ctx.GetRules(forceRefresh)
	if err != nil {
		return 0, err
	}
	for _, element := range frwRules {
		if element.Comment.GetValue() == comment {
			return element.Id.GetValue(), nil
		}
	}

	logging.CommonLog().Error("[DeviceController] FindRuleIdByComment Cannot find target rule")
	return 0, errors.New("cannot find target rule")
}

func (ctx *DeviceController) GetRulesFiltered(filter comment.IRuleComment, forceRefresh bool) ([]firewallCommon.FirewallRule, error) {
	frwRules, err := ctx.GetRules(forceRefresh)
	if err != nil {
		return frwRules, err
	}

	var filteredRules []firewallCommon.FirewallRule
	for _, element := range frwRules {
		if element.Comment.GetValue() == "" {
			continue
		}
		if filter.IsSameFamily(element.Comment.GetValue()) {
			filter.FromString(element.Comment.GetValue())
			filteredRules = append(filteredRules, element)
		}
	}
	return filteredRules, nil
}

func (ctx *DeviceController) CleanupTrashRules(comment comment.IRuleComment) error {
	/*
		Try to find rules for remove by:
		1) Rule prefix
		2) Controller name
		3) Endpoint hash
	*/
	var pendingRemove []uint64
	deviceResponse, err := ctx.Device.RunCommandWithReply(command.GetNew())
	if err != nil {
		logging.CommonLog().Error("[DeviceController] Error running command 'get rules' on device")
		return fmt.Errorf("error running command 'get rules' on device")
	}
	rules := deviceResponse.(*response.Get).GetRules()
	for _, rule := range rules {
		if comment.IsSameFamily(rule.Comment.GetValue()) {
			pendingRemove = append(pendingRemove, rule.Id.GetValue())
		}
	}
	return ctx.RemoveBatch(pendingRemove)
}

func (ctx *DeviceController) Start() error {
	ctx.deviceMutex.Lock()
	defer ctx.deviceMutex.Unlock()
	err := ctx.Device.Start()
	return err
}

func (ctx *DeviceController) Stop() error {
	ctx.deviceMutex.Lock()
	defer ctx.deviceMutex.Unlock()
	return ctx.Device.Stop()
}

func (ctx *DeviceController) RunCommandWithReply(cmd device.IDeviceCommand) (device.IDeviceResponse, error) {
	ctx.deviceMutex.Lock()
	defer ctx.deviceMutex.Unlock()
	switch cmd.GetType() {
	case device.DeviceCommandAdd:
	case device.DeviceCommandRemove:
		// Set needCacheUpdate to true after Add/Remove operations
		ctx.needCacheUpdate = true
	}
	return ctx.Device.RunCommandWithReply(cmd)
}
