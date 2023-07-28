package devices

import (
	"fmt"

	"github.com/altucor/http-knocker/device"
	"github.com/altucor/http-knocker/firewallProtocol"
	"github.com/altucor/http-knocker/logging"
	"gopkg.in/yaml.v3"

	"github.com/go-routeros/routeros"
)

type ConnectionRouterOsApi struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	Tls      bool   `yaml:"tls"`
}

type DeviceRouterOsApi struct {
	config ConnectionRouterOsApi
	client *routeros.Client
	tls    bool
}

func DeviceRouterOsNew(cfg ConnectionRouterOsApi) *DeviceRouterOsApi {
	ctx := &DeviceRouterOsApi{
		client: nil,
		config: ConnectionRouterOsApi{
			Username: cfg.Username,
			Password: cfg.Password,
			Host:     cfg.Host,
			Port:     cfg.Port,
			Tls:      cfg.Tls,
		},
	}
	return ctx
}

func DeviceRouterOsNewFromYaml(value *yaml.Node) (IDevice, error) {
	var cfg struct {
		Conn ConnectionRouterOsApi `yaml:"connection"`
	}
	if err := value.Decode(&cfg); err != nil {
		return nil, err
	}
	return DeviceRouterOsNew(cfg.Conn), nil
}

func (ctx *DeviceRouterOsApi) SetProtocol(protocol firewallProtocol.IFirewallProtocol) {
	// just do nothing, and keep object interface
	// we dont need protocol for this device
}

func (ctx *DeviceRouterOsApi) Start() error {
	logging.CommonLog().Info("[deviceMikrotik] Starting...")
	logging.CommonLog().Info("[deviceMikrotik] Starting... DONE")
	return nil
}

func (ctx *DeviceRouterOsApi) Stop() error {
	logging.CommonLog().Info("[deviceMikrotik] Stopping...")
	logging.CommonLog().Info("[deviceMikrotik] Stopping... DONE")
	return nil
}

func (ctx *DeviceRouterOsApi) Connect() {
	logging.CommonLog().Info("[deviceMikrotik] Connect called")
	var routerConnection *routeros.Client = nil
	var err error = nil
	if ctx.tls {
		/*
			routerConnection, err = routeros.DialTLS(
				ctx.device.Host+":"+fmt.Sprint(ctx.device.Port),
				ctx.device.Username,
				ctx.device.Password)
		*/
	} else {
		routerConnection, err = routeros.Dial(
			ctx.config.Host+":"+fmt.Sprint(ctx.config.Port),
			ctx.config.Username,
			ctx.config.Password)
	}

	if err != nil {
		logging.CommonLog().Fatal(err)
	}
	ctx.client = routerConnection
}

func (ctx *DeviceRouterOsApi) Disconnect() {
	logging.CommonLog().Info("[deviceMikrotik] Disconnect called")
	ctx.client.Close()
}

func (ctx *DeviceRouterOsApi) IsConnected() bool {
	return ctx.client != nil
}

func (ctx *DeviceRouterOsApi) RunCommandWithReply(command device.IDeviceCommand) (device.IDeviceResponse, error) {
	// response := ""
	// logging.CommonLog().Info("[deviceMikrotik] RunCommandWithReply called with =", cmd)
	// if ctx.IsConnected() {
	// 	reply, err := ctx.client.Run(strings.Fields(cmd)...)
	// 	if err != nil {
	// 		logging.CommonLog().Error(err)
	// 		return response, err
	// 	}
	// 	responseBytes, err := json.Marshal(*reply)
	// 	response = string(responseBytes)
	// 	if err != nil {
	// 		logging.CommonLog().Error(err)
	// 		return response, err
	// 	}
	// 	logging.CommonLog().Debug("[deviceMikrotik] RunCommandWithReply reply =", response)
	// } else {
	// 	logging.CommonLog().Error("[deviceMikrotik] Device is not connected, cannot run command")
	// }
	// return response, nil
	return nil, nil
}
