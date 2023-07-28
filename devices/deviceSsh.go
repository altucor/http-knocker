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
	protocol IFirewallSshProtocol
	client   *ssh.Client
	session  *ssh.Session
}

func DeviceSshNew(cfg ConnectionSSHCfg) *DeviceSsh {
	ctx := &DeviceSsh{
		client:   nil,
		config:   cfg,
		protocol: nil,
	}

	return ctx
}

func DeviceSshNewFromYaml(value *yaml.Node) (IDevice, error) {
	var cfg struct {
		Conn ConnectionSSHCfg `yaml:"connection"`
	}
	if err := value.Decode(&cfg); err != nil {
		return nil, err
	}
	return DeviceSshNew(cfg.Conn), nil
}

func (ctx *DeviceSsh) SetProtocol(protocol firewallProtocol.IFirewallProtocol) {
	ctx.protocol = protocol.(IFirewallSshProtocol)
}

func (ctx *DeviceSsh) Start() error {
	logging.CommonLog().Info("[deviceSsh Start] Starting...")
	ctx.ClientConnect()
	// ctx.SessionStart()
	logging.CommonLog().Info("[deviceSsh Start] Starting... DONE")
	return nil
}

func (ctx *DeviceSsh) Stop() error {
	logging.CommonLog().Info("[deviceSsh Stop] Stopping...")
	ctx.SessionStop()
	ctx.ClientDisconnect()
	logging.CommonLog().Info("[deviceSsh Stop] Stopping... DONE")
	return nil
}

func (ctx *DeviceSsh) ClientConnect() {
	logging.CommonLog().Info("[deviceSsh ClientConnect] Connect called")
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
		logging.CommonLog().Errorf("[deviceSsh ClientConnect] Connect error: %s\n", err)
	}
	ctx.client = sshConnection
}

func (ctx *DeviceSsh) ClientDisconnect() {
	logging.CommonLog().Info("[deviceSsh ClientDisconnect] Disconnect called")
	ctx.client.Close()
}

func (ctx *DeviceSsh) SessionStart() {
	var err error = nil
	ctx.session, err = ctx.client.NewSession()
	if err != nil {
		logging.CommonLog().Error("[deviceSsh SessionStart] RunSSHCommandWithReply NewSession error:", err)
	}
	// configure terminal mode
	modes := ssh.TerminalModes{
		ssh.ECHO: 0, // supress echo
	}
	// run terminal session
	if err := ctx.session.RequestPty("xterm", 50, 80, modes); err != nil {
		logging.CommonLog().Error("[deviceSsh SessionStart] RunSSHCommandWithReply RequestPty error:", err)
	}
}

func (ctx *DeviceSsh) SessionStop() {
	ctx.session.Close()
}

func (ctx *DeviceSsh) RunSSHCommandWithReply(cmd string) (string, error) {
	ctx.SessionStart()
	logging.CommonLog().Debug("[deviceSsh RunSSHCommandWithReply] Executing command: ", cmd)
	output, err := ctx.session.Output(cmd)
	if err != nil {
		logging.CommonLog().Error("[deviceSsh RunSSHCommandWithReply] Output error: ", err)
		logging.CommonLog().Error("[deviceSsh RunSSHCommandWithReply] Output: ", string(output))
		return "", err
	}
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
		logging.CommonLog().Error("[deviceSsh RunCommandWithReply] failed to convert cmd to IpTables: ", err)
		return &response.Add{}, err
	}
	output, err := ctx.RunSSHCommandWithReply(sshStr)
	if err != nil {
		logging.CommonLog().Error("[deviceSsh RunCommandWithReply] failed to execute command: ", err)
		return &response.Add{}, err
	}
	logging.CommonLog().Info("[deviceSsh RunCommandWithReply] reply = ", string(output))
	return ctx.protocol.From(output, command.GetType())
}
