package devices

import (
	"fmt"
	"net/http"

	"github.com/altucor/http-knocker/device"
	"github.com/altucor/http-knocker/firewallProtocol"
	"github.com/altucor/http-knocker/logging"
	"gopkg.in/yaml.v3"
)

type ConnectionRest struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Endpoint string `yaml:"endpoint"`
}

type IFirewallRestProtocol interface {
	firewallProtocol.IFirewallProtocol
	To(cmd device.IDeviceCommand, baseUrl string) (*http.Request, error)
	From(httpResponse *http.Response, cmdType device.DeviceCommandType) (device.IDeviceResponse, error)
}

type DeviceRest struct {
	config   ConnectionRest
	protocol IFirewallRestProtocol
}

func DeviceRestNew(cfg ConnectionRest, protocol IFirewallRestProtocol) *DeviceRest {
	ctx := &DeviceRest{
		config:   cfg,
		protocol: protocol,
	}
	return ctx
}

func DeviceRestNewFromYaml(value *yaml.Node, protocol IFirewallRestProtocol) (*DeviceRest, error) {
	var cfg struct {
		Conn ConnectionRest `yaml:"connection"`
	}
	if err := value.Decode(&cfg); err != nil {
		return nil, err
	}
	return DeviceRestNew(cfg.Conn, protocol), nil
}

func (ctx *DeviceRest) Start() error {
	logging.CommonLog().Info("[deviceRest] Starting...")
	logging.CommonLog().Info("[deviceRest] Starting... DONE")
	return nil
}

func (ctx *DeviceRest) Stop() error {
	logging.CommonLog().Info("[deviceRest] Stopping...")
	logging.CommonLog().Info("[deviceRest] Stopping... DONE")
	return nil
}

func (ctx *DeviceRest) isAvailable() bool {
	logging.CommonLog().Info("isAvailable called")
	return false
}

func (ctx *DeviceRest) executeRestCommand(request *http.Request) (*http.Response, error) {
	request.SetBasicAuth(ctx.config.Username, ctx.config.Password)
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		logging.CommonLog().Error("error making http request:", err)
		return nil, err
	}
	if res != nil {
		logging.CommonLog().Info("[deviceRest] status code:", res.StatusCode)
		logging.CommonLog().Debug("[deviceRest] response:", res)
	}
	return res, nil
}

func (ctx *DeviceRest) RunCommandWithReply(command device.IDeviceCommand) (device.IDeviceResponse, error) {
	var req *http.Request
	var err error = nil
	req, err = ctx.protocol.To(command, ctx.config.Endpoint)
	if err != nil {
		logging.CommonLog().Error("[deviceRest] RunCommandWithReply: Error marshaling command to REST:", err)
		return nil, fmt.Errorf("[deviceRest] RunCommandWithReply: Error marshaling command to REST: %s", err)
	}
	httpResponse, err := ctx.executeRestCommand(req)
	if err != nil {
		logging.CommonLog().Error("[deviceRest] RunCommandWithReply: Error executing command:", err)
		return nil, fmt.Errorf("[deviceRest] RunCommandWithReply: Error executing command: %s", err)
	}
	return ctx.protocol.From(httpResponse, command.GetType())
}
