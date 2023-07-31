package devices

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/altucor/http-knocker/device"
	"github.com/altucor/http-knocker/device/response"
	"github.com/altucor/http-knocker/firewallProtocol"
	"github.com/altucor/http-knocker/logging"
	"github.com/tarm/serial"
	"gopkg.in/yaml.v3"
)

type ConnectionSerialCfg struct {
	Name        string        `yaml:"name"`
	Baud        uint32        `yaml:"baud"`
	ReadTimeout time.Duration `yaml:"readTimeout"`
}

type outputCollector struct {
	buffer bytes.Buffer
}

type DeviceSerial struct {
	config          ConnectionSerialCfg
	protocol        IFirewallSshProtocol
	serialConfig    *serial.Config
	port            *serial.Port
	outputCollector outputCollector
}

func DeviceSerialNew(cfg ConnectionSerialCfg) *DeviceSerial {
	c := &serial.Config{
		Name:        cfg.Name,
		Baud:        int(cfg.Baud),
		ReadTimeout: cfg.ReadTimeout * time.Millisecond,
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
	go ctx.outputCollectorReader()
	logging.CommonLog().Info("[deviceSerial] Starting... DONE")
	return err
}

func (ctx *DeviceSerial) Stop() error {
	logging.CommonLog().Info("[deviceSerial] Stopping...")
	err := ctx.port.Close()
	logging.CommonLog().Info("[deviceSerial] Stopping... DONE")
	return err
}

func (ctx *DeviceSerial) outputCollectorReader() {
	read_data := make([]byte, 256)
	for {
		read_n, err := ctx.port.Read(read_data)
		if err != nil {
			logging.CommonLog().Errorf("[deviceSerial] outputCollectorReader Read_n %d Read error: %s", read_n, err)
			continue
		}
		written_n, err := ctx.outputCollector.buffer.Write(read_data[:read_n])
		if written_n != read_n {
			logging.CommonLog().Errorf("[deviceSerial] outputCollectorReader: not all data written to buffer Read_n %d, written_n %d Read error: %s", read_n, written_n, err)
			continue
		}
	}
}

func filterRules(input string, delimiter string) string {
	parts := strings.Split(input, delimiter)
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return parts[0]
}

func filterEscapeChars(in string) string {
	// https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797 -- Escape codes with description
	// \x1b[37;41;1m

	// To disable colors in MikroTik: https://forum.mikrotik.com/viewtopic.php?t=31350
	var reEscapeComamnds = regexp.MustCompile(`\x1b\[\d{0,};{0,}\d{0,};{0,}\d{0,}[a-zA-Z]`)
	return reEscapeComamnds.ReplaceAllString(in, ``)
}

func (ctx *DeviceSerial) readStringFromBuffer(n uint64) (string, error) {
	response := make([]byte, n)
	read_n, err := ctx.outputCollector.buffer.Read(response)
	if err != nil {
		logging.CommonLog().Errorf("[deviceSerial] readStringFromBuffer Read_n %d Read error: %s", read_n, err)
		return "", err
	}

	return string(response[:read_n]), nil
}

func (ctx *DeviceSerial) getResponseForCmd(cmd string, expectedOutputSize uint64) (string, error) {
	str, err := ctx.readStringFromBuffer(expectedOutputSize)
	if err != nil {
		return str, err
	}
	// str = strings.ReplaceAll(str, "\r", "")
	str = filterEscapeChars(str)
	str = filterRules(str, cmd)
	return str, nil
}

func (ctx *DeviceSerial) RunSerialCommandWithReply(cmd string, timeout time.Duration) (string, error) {
	cmd += "\r"
	written_n, err := ctx.port.Write([]byte(cmd))
	if err != nil {
		logging.CommonLog().Errorf("[deviceSerial] RunSerialCommandWithReply Written_n %d Write error: %s", written_n, err)
		return "", err
	}
	time.Sleep(timeout)
	return ctx.getResponseForCmd(cmd, 8*1024*1024)
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
	output, err := ctx.RunSerialCommandWithReply(serialStr, time.Second)
	if err != nil {
		logging.CommonLog().Error("[deviceSerial] RunCommandWithReply failed to execute command: ", err)
		return &response.Add{}, err
	}
	logging.CommonLog().Info("[deviceSerial] RunCommandWithReply reply = ", string(output))
	return ctx.protocol.From(output, command.GetType())
}
