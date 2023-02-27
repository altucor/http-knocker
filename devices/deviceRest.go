package devices

import (
	"bytes"
	"errors"
	"fmt"
	"http-knocker/deviceCommon"
	"http-knocker/deviceResponse"
	"http-knocker/logging"
	"net/http"

	"golang.org/x/exp/slices"
)

type ConnectionRest struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Enpoint  string `yaml:"endpoint"`
	Tls      bool   `yaml:"tls"`
}

type DeviceRest struct {
	config             ConnectionRest
	supportedProtocols []DeviceProtocol
}

func DeviceRestNew(cfg ConnectionRest) *DeviceRest {
	ctx := &DeviceRest{
		config: cfg,
		supportedProtocols: []DeviceProtocol{
			PROTOCOL_ROUTER_OS_REST,
		},
	}
	return ctx
}

func (ctx *DeviceRest) Start() error {
	return nil
}

func (ctx *DeviceRest) Stop() error {
	return nil
}

func (ctx *DeviceRest) isAvailable() bool {
	logging.CommonLog().Info("isAvailable called")
	return false
}

func (ctx *DeviceRest) GetSupportedProtocols() []DeviceProtocol {
	return ctx.supportedProtocols
}

func (ctx *DeviceRest) GetType() DeviceType {
	return DeviceTypeRest
}

func (ctx *DeviceRest) executeRestCommand(method string, url string, body string) (http.Response, error) {
	req, err := http.NewRequest(method, ctx.config.Enpoint+url, bytes.NewReader([]byte(body)))
	if err != nil {
		logging.CommonLog().Error("could not create request: %s\n", err)
	}
	req.SetBasicAuth(ctx.config.Username, ctx.config.Password)
	req.Header.Set("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logging.CommonLog().Error("error making http request: %s\n", err)
		return http.Response{}, err
	}
	if res != nil {
		logging.CommonLog().Info("[deviceRest] status code:", res.StatusCode)
		logging.CommonLog().Debug("[deviceRest] response:", res)
	}
	return *res, nil
}

func (ctx *DeviceRest) RunCommandWithReply(command deviceCommon.IDeviceCommand, proto DeviceProtocol) (deviceCommon.IDeviceResponse, error) {
	if !slices.Contains(ctx.supportedProtocols, proto) {
		return nil, errors.New(fmt.Sprintf("[deviceRest] RunCommandWithReply: Error not supported protocol: %s", proto))
	}
	var method string = ""
	var url string = ""
	var body string = ""
	var err error = nil
	switch proto {
	case PROTOCOL_ROUTER_OS_REST:
		method, url, body, err = command.Rest()
	}
	if err != nil {
		logging.CommonLog().Error("[deviceRest] RunCommandWithReply: Error marshaling command to REST %s", err)
		return nil, errors.New(fmt.Sprintf("[deviceRest] RunCommandWithReply: Error marshaling command to REST %s", err))
	}
	httpResponse, err := ctx.executeRestCommand(method, url, body)
	if err != nil {
		logging.CommonLog().Error("[deviceRest] RunCommandWithReply: Error executing command %s", err)
		return nil, errors.New(fmt.Sprintf("[deviceRest] RunCommandWithReply: Error executing command %s", err))
	}
	switch command.GetType() {
	case deviceCommon.DeviceCommandAdd:
		return deviceResponse.AddFromRouterOsRest(httpResponse)
	case deviceCommon.DeviceCommandGet:
		return deviceResponse.GetFromRouterOsRest(httpResponse)
	case deviceCommon.DeviceCommandRemove:
		return deviceResponse.RemoveFromRouterOsRest(httpResponse)
	case deviceCommon.DeviceCommandMove:
		return deviceResponse.MoveFromRouterOsRest(httpResponse)
	default:
		logging.CommonLog().Error("[deviceRest] Unknown response type")
		return nil, errors.New("[deviceRest] Unknown response type")
	}
}

// func (ctx *DeviceRest) RunCommand(command deviceCommon.IDeviceCommand) error {
// 	_, err := ctx.RunCommandWithReply(command)
// 	return err
// }
