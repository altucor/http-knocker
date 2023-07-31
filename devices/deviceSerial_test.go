package devices

import (
	"fmt"
	"testing"
	"time"

	"github.com/altucor/http-knocker/device/command"
	"github.com/altucor/http-knocker/firewallProtocol"
)

func Test_DeviceSerial_ConnectDisconnect(t *testing.T) {
	c := ConnectionSerialCfg{
		Name:        "/dev/tty.usbserial-21230",
		Baud:        115200,
		ReadTimeout: 500,
	}
	dev := DeviceSerialNew(c)
	dev.Start()
	reply, err := dev.RunSerialCommandWithReply("?", time.Second)
	dev.Stop()
	if err != nil {
		t.Error("Error from dev.RunSerialCommandWithReply:", err)
	}

	// \r\n\r\n\x1b[m\x1b[35mbeep\x1b[m\x1b[33m -- \r\n\x1b[\x00\x00
	// x1b[9999B[admin@1gb-switch] > \x1b[K\r\n\r\r\r\r

	t.Log("\n")
	if len(reply) == 0 {
		t.Errorf("\"DeviceSerial_ConnectDisconnect\" FAILED")
	} else {
		t.Logf("\"DeviceSerial_ConnectDisconnect\" SUCCEDED")
	}
}

func Test_DeviceSerial_GetFirewallRules(t *testing.T) {
	c := ConnectionSerialCfg{
		Name:        "/dev/tty.usbserial-21230",
		Baud:        115200,
		ReadTimeout: 500,
	}
	dev := DeviceSerialNew(c)
	dev.Start()
	reply, err := dev.RunSerialCommandWithReply(" /ip firewall filter print without-paging", time.Second)
	dev.Stop()
	if err != nil {
		t.Error("Error getting firewall rules:", err)
	}
	fmt.Print(reply)

	t.Log("\n")
	if len(reply) == 0 {
		t.Errorf("\"DeviceSerial_ConnectDisconnect\" FAILED")
	} else {
		t.Logf("\"DeviceSerial_ConnectDisconnect\" SUCCEDED")
	}
}

func Test_DeviceSerial_RequestFirewallRulesWithParsing(t *testing.T) {
	c := ConnectionSerialCfg{
		Name:        "/dev/tty.usbserial-21230",
		Baud:        115200,
		ReadTimeout: 500,
	}
	dev := DeviceSerialNew(c)
	dev.Start()
	dev.SetProtocol(firewallProtocol.GetProtocolStorage().GetProtocolByName("terminal-router-os"))
	reply, err := dev.RunCommandWithReply(command.GetNew())
	dev.Stop()
	if err != nil {
		t.Error("Error getting firewall rules:", err)
	}
	fmt.Print(reply)

	t.Log("\n")
	// if len(reply) == 0 {
	// 	t.Errorf("\"DeviceSerial_ConnectDisconnect\" FAILED")
	// } else {
	// 	t.Logf("\"DeviceSerial_ConnectDisconnect\" SUCCEDED")
	// }
}
