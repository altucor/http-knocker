package devices

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/altucor/http-knocker/device"
	"github.com/altucor/http-knocker/firewallProtocol"
	"github.com/altucor/http-knocker/logging"
	"gopkg.in/yaml.v3"
)

type ConnectionRest struct {
	Username           string `yaml:"username"`
	Password           string `yaml:"password"`
	Endpoint           string `yaml:"endpoint"`
	InsecureSkipVerify bool   `yaml:"insecure-skip-verify"`
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

func DeviceRestNew(cfg ConnectionRest) *DeviceRest {
	ctx := &DeviceRest{
		config:   cfg,
		protocol: nil,
	}
	return ctx
}

func DeviceRestNewFromYaml(value *yaml.Node) (IDevice, error) {
	var cfg struct {
		Conn ConnectionRest `yaml:"connection"`
	}
	if err := value.Decode(&cfg); err != nil {
		return nil, err
	}
	return DeviceRestNew(cfg.Conn), nil
}

func (ctx *DeviceRest) SetProtocol(protocol firewallProtocol.IFirewallProtocol) {
	ctx.protocol = protocol.(IFirewallRestProtocol)
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
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: ctx.config.InsecureSkipVerify},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Do(request)
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
