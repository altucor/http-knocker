package devices

import (
	"github.com/altucor/http-knocker/comment"
	"github.com/altucor/http-knocker/device"
	"github.com/altucor/http-knocker/device/command"
	"github.com/altucor/http-knocker/device/response"
	"github.com/altucor/http-knocker/firewallProtocol"
	"github.com/altucor/http-knocker/logging"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Type           string `yaml:"type"`
	Protocol       string `yaml:"protocol"`
	CleanupOnStart bool   `yaml:"cleanup-on-start"`
}

type DeviceController struct {
	Device IDevice
	Config Config
}

func DeviceControllerNew(device IDevice) *DeviceController {
	ctx := DeviceController{
		Device: nil,
		Config: Config{
			Type:           "",
			Protocol:       "",
			CleanupOnStart: false,
		},
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
	device, err := deviceStorage.GetDeviceConstructor(ctx.Config.Type).(func(*yaml.Node, firewallProtocol.IFirewallProtocol) (IDevice, error))(value, protocol)
	if err != nil {
		return err
	}
	ctx.Device = device
	return nil
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
	deviceResponse, err := ctx.Device.RunCommandWithReply(command.GetNew())
	if err != nil {
		logging.CommonLog().Error("[DeviceController] Error running command 'get rules' on device")
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
	err := ctx.Device.Start()
	// if ctx.Config.CleanupOnStart {
	// 	// ctx.CleanupTrashRules()
	// }
	return err
}

func (ctx *DeviceController) Stop() error {
	return ctx.Device.Stop()
}

func (ctx *DeviceController) RunCommandWithReply(cmd device.IDeviceCommand) (device.IDeviceResponse, error) {
	return ctx.Device.RunCommandWithReply(cmd)
}
