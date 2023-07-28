package devices

import (
	"errors"
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
	Type           string `yaml:"type"`
	Protocol       string `yaml:"protocol"`
	CleanupOnStart bool   `yaml:"cleanup-on-start"`
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
			Type:           "",
			Protocol:       "",
			CleanupOnStart: false,
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

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func (ctx *DeviceController) CleanupTrashRules(knownIdentifiers []string) {
	/*
		Try to find rules for remove by:
		1) Rule prefix
		2) Controller name
		3) Endpoint hash
	*/
	deviceResponse, err := ctx.Device.RunCommandWithReply(command.GetNew())
	if err != nil {
		logging.CommonLog().Error("[DeviceController] Error running command 'get rules' on device")
		return
	}
	rules := deviceResponse.(*response.Get).GetRules()
	for _, rule := range rules {
		commentData, err := comment.BasicCommentNewFromString(rule.Comment.GetString(), "-")
		if err != nil {
			logging.CommonLog().Error("[DeviceController] Error processing firewall rule comment", err)
			continue
		}
		if commentData.GetEndpointHash() == "" {
			continue
		}
		if !contains(knownIdentifiers, commentData.GetEndpointHash()) {
			_, err := ctx.Device.RunCommandWithReply(command.RemoveNew(rule.Id.GetValue()))
			if err != nil {
				logging.CommonLog().Error("[DeviceController] Error removing unrelated rule:", rule)
			}
		}
	}
}

func (ctx *DeviceController) Start() error {
	ctx.deviceMutex.Lock()
	defer ctx.deviceMutex.Unlock()
	err := ctx.Device.Start()
	// if ctx.Config.CleanupOnStart {
	// 	// ctx.CleanupTrashRules()
	// }
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
