package devices

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"sync"
	"time"

	"github.com/altucor/http-knocker/common"
	"github.com/altucor/http-knocker/device"
	"github.com/altucor/http-knocker/device/response"
	"github.com/altucor/http-knocker/firewallProtocol"
	"github.com/altucor/http-knocker/logging"
	"gopkg.in/yaml.v3"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type ConnectionSSHCfg struct {
	Username                   string                 `yaml:"username"`
	Password                   string                 `yaml:"password"`
	Host                       string                 `yaml:"host"`
	Port                       uint16                 `yaml:"port"`
	KnownHosts                 string                 `yaml:"knownHosts"`
	PrivateKeyPath             string                 `yaml:"private-key"`
	PrivateKeyPassphrase       string                 `yaml:"private-key-passphrase"`
	ConnectionKeepAliveTimeout common.DurationSeconds `yaml:"connection-keep-alive"`
}

type IFirewallSshProtocol interface {
	firewallProtocol.IFirewallProtocol
	To(cmd device.IDeviceCommand) (string, error)
	From(responseData string, cmdType device.DeviceCommandType) (device.IDeviceResponse, error)
}

type DeviceSsh struct {
	config                   ConnectionSSHCfg
	protocol                 IFirewallSshProtocol
	client                   *ssh.Client
	muClient                 sync.Mutex
	connectionKeepAliveTimer time.Timer
	lastCLientCloseError     error
	session                  *ssh.Session
}

func DeviceSshNew(cfg ConnectionSSHCfg) *DeviceSsh {
	ctx := &DeviceSsh{
		config:                   cfg,
		protocol:                 nil,
		client:                   nil,
		lastCLientCloseError:     nil,
		connectionKeepAliveTimer: *time.NewTimer(cfg.ConnectionKeepAliveTimeout.GetValue()),
		session:                  nil,
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

func (ctx *DeviceSsh) clientClose() error {
	ctx.muClient.Lock()
	defer ctx.muClient.Unlock()
	ctx.connectionKeepAliveTimer.Stop()
	if ctx.client != nil {
		ctx.lastCLientCloseError = ctx.client.Close()
	}
	ctx.client = nil
	return ctx.lastCLientCloseError
}

func (ctx *DeviceSsh) keepAliveHandler() {
	for {
		<-ctx.connectionKeepAliveTimer.C
		ctx.clientClose()
		logging.CommonLog().Error("[deviceSsh keepAliveHandler] Tiemout reached Close last error:", ctx.lastCLientCloseError)
	}
}

func (ctx *DeviceSsh) clientConnectionMonitor() {
	for {
		waitErr := ctx.client.Wait()
		logging.CommonLog().Error("[deviceSsh clientConnectionMonitor] Wait last error:", waitErr)
		ctx.clientClose()
		logging.CommonLog().Error("[deviceSsh clientConnectionMonitor] Close last error:", ctx.lastCLientCloseError)
		break
	}
}

func (ctx *DeviceSsh) hostKeyCallback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	logging.CommonLog().Info("[deviceSsh hostKeyCallback] called:", hostname, remote, key)
	return nil
}

func (ctx *DeviceSsh) getAuthMethod() (ssh.AuthMethod, error) {
	var err error = nil
	if ctx.config.PrivateKeyPath != "" {
		var key []byte
		key, err = ioutil.ReadFile(ctx.config.PrivateKeyPath)
		if err != nil {
			return nil, err
		}
		var signer ssh.Signer = nil
		if ctx.config.PrivateKeyPassphrase != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(ctx.config.PrivateKeyPassphrase))
		} else {
			signer, err = ssh.ParsePrivateKey(key)
		}
		if err != nil {
			return nil, err
		}
		return ssh.PublicKeys(signer), nil
	}
	if ctx.config.Password != "" {
		return ssh.Password(ctx.config.Password), nil
	}

	return nil, errors.New("cannot prepare auth method")
}

func (ctx *DeviceSsh) clientConnect() error {
	var err error = nil
	authMethod, err := ctx.getAuthMethod()
	if err != nil {
		logging.CommonLog().Error("[deviceSsh clientConnect] error getting auth method:", err)
		return err
	}

	var hostKeyCallback ssh.HostKeyCallback = ctx.hostKeyCallback
	if ctx.config.KnownHosts != "" {
		hostKeyCallback, err = knownhosts.New(ctx.config.KnownHosts)
		if err != nil {
			logging.CommonLog().Error("[deviceSsh clientConnect] error parsing KnownHosts:", err)
			return err
		}
	}

	config := &ssh.ClientConfig{
		User: ctx.config.Username,
		Auth: []ssh.AuthMethod{
			authMethod,
		},
		HostKeyCallback: hostKeyCallback,
		Timeout:         time.Second * 10,
	}
	ctx.client, err = ssh.Dial("tcp", ctx.config.Host+":"+fmt.Sprint(ctx.config.Port), config)
	if err == nil {
		go ctx.clientConnectionMonitor()
		ctx.connectionKeepAliveTimer.Reset(ctx.config.ConnectionKeepAliveTimeout.GetValue())
		go ctx.keepAliveHandler()
	}
	return err
}

func (ctx *DeviceSsh) Start() error {
	logging.CommonLog().Info("[deviceSsh Start] Starting...")
	ctx.muClient.Lock()
	defer ctx.muClient.Unlock()
	err := ctx.clientConnect()
	logging.CommonLog().Info("[deviceSsh Start] Starting... DONE")
	return err
}

func (ctx *DeviceSsh) Stop() error {
	logging.CommonLog().Info("[deviceSsh Stop] Stopping...")
	err := ctx.clientClose()
	logging.CommonLog().Info("[deviceSsh Stop] Stopping... DONE")
	return err
}

func (ctx *DeviceSsh) sessionStart() error {
	var err error = nil
	if ctx.client == nil {
		logging.CommonLog().Error("[deviceSsh sessionStart] Error got nil client")
		return errors.New("error got nil client")
	}
	ctx.session, err = ctx.client.NewSession()
	if err != nil {
		return err
	}
	// configure terminal mode
	modes := ssh.TerminalModes{
		ssh.ECHO: 0, // supress echo
	}
	// run terminal session
	err = ctx.session.RequestPty("xterm", 50, 80, modes)
	return err
}

func (ctx *DeviceSsh) sessionStop() {
	ctx.session.Close()
}

func (ctx *DeviceSsh) RunSSHCommandWithReply(cmd string) (string, error) {
	ctx.muClient.Lock()
	defer ctx.muClient.Unlock()
	if ctx.client == nil {
		if err := ctx.clientConnect(); err != nil {
			logging.CommonLog().Error("[deviceSsh RunSSHCommandWithReply] Error connecting client: ", err)
		}
	}
	if err := ctx.sessionStart(); err != nil {
		logging.CommonLog().Error("[deviceSsh RunSSHCommandWithReply] Error starting session: ", err)
	}
	defer ctx.sessionStop()
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
