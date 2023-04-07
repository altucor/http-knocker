package devices

import (
	"fmt"

	"github.com/altucor/http-knocker/device"
	"github.com/altucor/http-knocker/device/response"
	"github.com/altucor/http-knocker/firewallProtocol"
	"github.com/altucor/http-knocker/logging"
	"gopkg.in/yaml.v3"

	"golang.org/x/crypto/ssh"
)

type ConnectionSSHCfg struct {
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Host       string `yaml:"host"`
	Port       uint16 `yaml:"port"`
	KnownHosts string `yaml:"knownHosts"`
}

type IFirewallSshProtocol interface {
	firewallProtocol.IFirewallProtocol
	To(cmd device.IDeviceCommand) (string, error)
	From(responseData string, cmdType device.DeviceCommandType) (device.IDeviceResponse, error)
}

type DeviceSsh struct {
	config   ConnectionSSHCfg
	client   *ssh.Client
	protocol IFirewallSshProtocol
}

func DeviceSshNew(cfg ConnectionSSHCfg, protocol IFirewallSshProtocol) *DeviceSsh {
	ctx := &DeviceSsh{
		client:   nil,
		config:   cfg,
		protocol: protocol,
	}

	return ctx
}

func DeviceSshNewFromYaml(value *yaml.Node, protocol IFirewallSshProtocol) (*DeviceSsh, error) {
	var cfg struct {
		Conn ConnectionSSHCfg `yaml:"connection"`
	}
	if err := value.Decode(&cfg); err != nil {
		return nil, err
	}
	return DeviceSshNew(cfg.Conn, protocol), nil
}

func (ctx *DeviceSsh) Start() error {
	logging.CommonLog().Info("[deviceSsh] Starting...")
	logging.CommonLog().Info("[deviceSsh] Starting... DONE")
	return nil
}

func (ctx *DeviceSsh) Stop() error {
	logging.CommonLog().Info("[deviceSsh] Stopping...")
	logging.CommonLog().Info("[deviceSsh] Stopping... DONE")
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

func (ctx *DeviceSsh) RunCommandWithReply(command device.IDeviceCommand) (device.IDeviceResponse, error) {
	var sshStr string = ""
	var err error = nil

	if ctx.protocol == nil {
		return nil, fmt.Errorf("protocol is not set")
	}

	sshStr, err = ctx.protocol.To(command)
	if err != nil {
		logging.CommonLog().Error("[deviceSsh] RunCommandWithReply failed to convert cmd to IpTables: %s\n", err)
		return &response.Add{}, err
	}
	output, err := ctx.RunSSHCommandWithReply(sshStr)
	if err != nil {
		logging.CommonLog().Error("[deviceSsh] RunCommandWithReply failed to execute command: %s\n", err)
		return &response.Add{}, err
	}
	logging.CommonLog().Info("[deviceSsh] RunCommand reply =", string(output))
	return ctx.protocol.From(output, command.GetType())
}
