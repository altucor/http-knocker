package devices

import (
	"fmt"

	"github.com/altucor/http-knocker/device"
	"github.com/altucor/http-knocker/device/response"
	"github.com/altucor/http-knocker/firewallProtocol"
	"github.com/altucor/http-knocker/logging"
	"github.com/tarm/serial"
	"gopkg.in/yaml.v3"
)

type ConnectionSerialCfg struct {
	Name string `yaml:"name"`
	Baud uint32 `yaml:"baud"`
}

type DeviceSerial struct {
	config       ConnectionSerialCfg
	protocol     IFirewallSshProtocol
	serialConfig *serial.Config
	port         *serial.Port
}

func DeviceSerialNew(cfg ConnectionSerialCfg) *DeviceSerial {
	c := &serial.Config{
		Name: cfg.Name,
		Baud: int(cfg.Baud),
	}
	serial.OpenPort(c)
	ctx := &DeviceSerial{
		config:       cfg,
		serialConfig: c,
		port:         nil,
	}
	return ctx
}

func DeviceSerialNewFromYaml(value *yaml.Node) (IDevice, error) {
	var cfg struct {
		Conn ConnectionSerialCfg `yaml:"connection"`
	}
	if err := value.Decode(&cfg); err != nil {
		return nil, err
	}
	return DeviceSerialNew(cfg.Conn), nil
}

func (ctx *DeviceSerial) SetProtocol(protocol firewallProtocol.IFirewallProtocol) {
	ctx.protocol = protocol.(IFirewallSshProtocol)
}

func (ctx *DeviceSerial) Start() error {
	logging.CommonLog().Info("[deviceSerial] Starting...")
	port, err := serial.OpenPort(ctx.serialConfig)
	if err != nil {
		logging.CommonLog().Error("[deviceSerial] error opening port:", ctx.serialConfig.Name)
	}
	ctx.port = port
	logging.CommonLog().Info("[deviceSerial] Starting... DONE")
	return err
}

func (ctx *DeviceSerial) Stop() error {
	logging.CommonLog().Info("[deviceSerial] Stopping...")
	err := ctx.port.Close()
	logging.CommonLog().Info("[deviceSerial] Stopping... DONE")
	return err
}

func (ctx *DeviceSerial) RunSerialCommandWithReply(cmd string) (string, error) {

	written_n, err := ctx.port.Write([]byte(cmd))
	if err != nil {
		logging.CommonLog().Errorf("[deviceSerial] RunSerialCommandWithReply Written_n %d Write error: %s", written_n, err)
		return "", err
	}
	var response []byte
	read_n, err := ctx.port.Read(response)
	if err != nil {
		logging.CommonLog().Errorf("[deviceSerial] RunSerialCommandWithReply Read_n %d Read error: %s", read_n, err)
		return "", err
	}
	return string(response), nil
}

func (ctx *DeviceSerial) RunCommandWithReply(command device.IDeviceCommand) (device.IDeviceResponse, error) {
	var serialStr string = ""
	var err error = nil

	if ctx.protocol == nil {
		return nil, fmt.Errorf("protocol is not set")
	}

	serialStr, err = ctx.protocol.To(command)
	if err != nil {
		logging.CommonLog().Error("[deviceSerial] RunCommandWithReply failed to convert cmd to IpTables: ", err)
		return &response.Add{}, err
	}
	output, err := ctx.RunSerialCommandWithReply(serialStr)
	if err != nil {
		logging.CommonLog().Error("[deviceSerial] RunCommandWithReply failed to execute command: ", err)
		return &response.Add{}, err
	}
	logging.CommonLog().Info("[deviceSerial] RunCommandWithReply reply = ", string(output))
	return ctx.protocol.From(output, command.GetType())
}
