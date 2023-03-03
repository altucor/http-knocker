package devices

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/altucor/http-knocker/logging"

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

func DeviceRouterOsNew(cfg DeviceConnectionDesc) *DeviceRouterOsApi {
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

func (ctx *DeviceRouterOsApi) Start() error {
	return nil
}

func (ctx *DeviceRouterOsApi) Stop() error {
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

func (ctx *DeviceRouterOsApi) GetType() DeviceType {
	return DeviceTypeRouterOs
}

func (ctx *DeviceRouterOsApi) IsConnected() bool {
	return ctx.client != nil
}

func (ctx *DeviceRouterOsApi) RunCommandWithReply(cmd string) (string, error) {
	response := ""
	logging.CommonLog().Info("[deviceMikrotik] RunCommandWithReply called with =", cmd)
	if ctx.IsConnected() {
		reply, err := ctx.client.Run(strings.Fields(cmd)...)
		if err != nil {
			logging.CommonLog().Error(err)
			return response, err
		}
		responseBytes, err := json.Marshal(*reply)
		response = string(responseBytes)
		if err != nil {
			logging.CommonLog().Error(err)
			return response, err
		}
		logging.CommonLog().Debug("[deviceMikrotik] RunCommandWithReply reply =", response)
	} else {
		logging.CommonLog().Error("[deviceMikrotik] Device is not connected, cannot run command")
	}
	return response, nil
}

func (ctx *DeviceRouterOsApi) RunCommand(cmd string) error {
	_, err := ctx.RunCommandWithReply(cmd)
	return err
}
