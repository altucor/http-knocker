package devices

import (
	"errors"
	"fmt"

	"github.com/altucor/http-knocker/deviceCommon"
	"github.com/altucor/http-knocker/deviceResponse"
	"github.com/altucor/http-knocker/logging"

	"golang.org/x/crypto/ssh"
	"golang.org/x/exp/slices"
)

type ConnectionSSHCfg struct {
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Host       string `yaml:"host"`
	Port       uint16 `yaml:"port"`
	KnownHosts string `yaml:"knownHosts"`
}

type DeviceSsh struct {
	config             ConnectionSSHCfg
	supportedProtocols []DeviceProtocol
	client             *ssh.Client
}

func DeviceSshNew(cfg ConnectionSSHCfg) *DeviceSsh {
	logging.CommonLog().Debug(cfg.Host + ":" + fmt.Sprint(cfg.Port))
	ctx := &DeviceSsh{
		client: nil,
		supportedProtocols: []DeviceProtocol{
			PROTOCOL_IP_TABLES,
		},
		config: cfg,
	}

	return ctx
}

func (ctx *DeviceSsh) Start() error {
	return nil
}

func (ctx *DeviceSsh) Stop() error {
	return nil
}

func (ctx *DeviceSsh) Connect() {
	logging.CommonLog().Info("[deviceSsh] Connect called")
	//hostKeyCallback, err := knownhosts.New("/home/debian11/.ssh/known_hosts")
	//if err != nil {
	// logging.CommonLog().Fatal(err)
	//}
	config := &ssh.ClientConfig{
		User: ctx.config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(ctx.config.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshConnection, err := ssh.Dial("tcp", ctx.config.Host+":"+fmt.Sprint(ctx.config.Port), config)
	if err != nil {
		logging.CommonLog().Error("[deviceSsh] Connect error: %s\n", err)
	}
	ctx.client = sshConnection
}

func (ctx *DeviceSsh) Disconnect() {
	logging.CommonLog().Info("[deviceSsh] Disconnect called")
	ctx.client.Close()
}

func (ctx *DeviceSsh) GetSupportedProtocols() []DeviceProtocol {
	return ctx.supportedProtocols
}

func (ctx *DeviceSsh) GetType() DeviceType {
	return DeviceTypeSsh
}

func (ctx *DeviceSsh) RunSSHCommandWithReply(cmd string) (string, error) {
	ctx.Connect()
	defer ctx.Disconnect()
	logging.CommonLog().Info("[deviceSsh] RunSSHCommandWithReply called with =", cmd)
	session, err := ctx.client.NewSession()
	if err != nil {
		logging.CommonLog().Error("[deviceSsh] RunSSHCommandWithReply NewSession error:", err)
		return "", err
	}
	defer session.Close()

	// configure terminal mode
	modes := ssh.TerminalModes{
		ssh.ECHO: 0, // supress echo
	}
	// run terminal session
	if err := session.RequestPty("xterm", 50, 80, modes); err != nil {
		logging.CommonLog().Error("[deviceSsh] RunSSHCommandWithReply RequestPty error:", err)
		return "", err
	}
	// start remote shell
	output, err := session.Output(cmd)
	if err != nil {
		logging.CommonLog().Error("[deviceSsh] RunSSHCommandWithReply Output error:", err)
		return "", err
	}
	logging.CommonLog().Debug("[deviceSsh] RunSSHCommandWithReply reply =", string(output))
	return string(output), nil
}

func (ctx *DeviceSsh) RunCommandWithReply(command deviceCommon.IDeviceCommand, proto DeviceProtocol) (deviceCommon.IDeviceResponse, error) {
	if !slices.Contains(ctx.supportedProtocols, proto) {
		return nil, errors.New(fmt.Sprintf("[deviceSsh] RunCommandWithReply: Error not supported protocol: %s", proto))
	}
	var ipTablesStr string = ""
	var err error = nil
	switch proto {
	case PROTOCOL_IP_TABLES:
		ipTablesStr, err = command.IpTables()
	}
	if err != nil {
		logging.CommonLog().Error("[deviceSsh] RunCommandWithReply failed to convert cmd to IpTables: %s\n", err)
		return deviceResponse.Add{}, err
	}
	output, err := ctx.RunSSHCommandWithReply(ipTablesStr)
	if err != nil {
		logging.CommonLog().Error("[deviceSsh] RunCommandWithReply failed to execute command: %s\n", err)
		return deviceResponse.Add{}, err
	}
	logging.CommonLog().Info("[deviceSsh] RunCommand reply =", string(output))
	switch command.GetType() {
	case deviceCommon.DeviceCommandAdd:
		return deviceResponse.AddFromIpTables(output)
	case deviceCommon.DeviceCommandGet:
		return deviceResponse.GetFromIpTables(output)
	case deviceCommon.DeviceCommandRemove:
		return deviceResponse.RemoveFromIpTables(output)
	case deviceCommon.DeviceCommandMove:
		return deviceResponse.MoveFromIpTables(output)
	default:
		return nil, errors.New("[deviceSsh] Unknown response type")
	}
}
