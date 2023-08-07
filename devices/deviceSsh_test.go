package devices

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/altucor/http-knocker/logging"
	"gopkg.in/yaml.v3"
)

func deviceSshFromYaml() ConnectionSSHCfg {
	cfg := ConnectionSSHCfg{}
	bytes, err := ioutil.ReadFile("sshDevice.yaml")
	if err != nil {
		logging.CommonLog().Errorf("[deviceFromYaml] Error reading file: %s", err)
		return cfg
	}
	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		logging.CommonLog().Errorf("[deviceFromYaml] Error unmarshaling yaml file: %s", err)
		return cfg
	}
	return cfg
}

func Test_DeviceSsh_Connect_Disconnect(t *testing.T) {
	input := "ls /"
	expected := "bin   dev  home  lib32	libx32	mnt  proc  run	 srv  tmp  var\r\nboot  etc  lib	 lib64	media	opt  root  sbin  sys  usr\r\n"
	var err error = nil

	ssh := DeviceSshNew(deviceSshFromYaml())
	err = ssh.Start()
	if err != nil {
		t.Error("Error Starting device:", err)
	}
	result, err := ssh.RunSSHCommandWithReply(input)
	if err != nil {
		t.Error("Error executing command:", err)
	}
	err = ssh.Stop()
	if err != nil {
		t.Error("Error Stopping device:", err)
	}

	if result != expected {
		t.Errorf("\"DeviceSsh_Connect_Disconnect('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
	} else {
		t.Logf("\"DeviceSsh_Connect_Disconnect('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
	}
}

func Test_DeviceSsh_Reconnect(t *testing.T) {
	input := "ls /"
	expected := "bin   dev  home  lib32	libx32	mnt  proc  run	 srv  tmp  var\r\nboot  etc  lib	 lib64	media	opt  root  sbin  sys  usr\r\n"
	var err error = nil

	ssh := DeviceSshNew(deviceSshFromYaml())

	// First round of connect
	err = ssh.Start()
	if err != nil {
		t.Error("Error Starting device:", err)
	}
	result, err := ssh.RunSSHCommandWithReply(input)
	if err != nil {
		t.Error("Error executing command:", err)
	}
	err = ssh.Stop()
	if err != nil {
		t.Error("Error Stopping device:", err)
	}
	if result != expected {
		t.Errorf("\"DeviceSsh_Reconnect Round 1 ('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
	} else {
		t.Logf("\"DeviceSsh_Reconnect Round 1 ('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
	}

	// Second round of connect
	err = ssh.Start()
	if err != nil {
		t.Error("Error Starting device:", err)
	}
	result, err = ssh.RunSSHCommandWithReply(input)
	if err != nil {
		t.Error("Error executing command:", err)
	}
	err = ssh.Stop()
	if err != nil {
		t.Error("Error Stopping device:", err)
	}
	if result != expected {
		t.Errorf("\"DeviceSsh_Reconnect Round 2 ('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
	} else {
		t.Logf("\"DeviceSsh_Reconnect Round 2 ('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
	}
}

func Test_DeviceSsh_Reconnect_WithTimeout(t *testing.T) {
	input := "ls /"
	expected := "bin   dev  home  lib32	libx32	mnt  proc  run	 srv  tmp  var\r\nboot  etc  lib	 lib64	media	opt  root  sbin  sys  usr\r\n"
	var err error = nil

	ssh := DeviceSshNew(deviceSshFromYaml())

	// First round of connect
	err = ssh.Start()
	if err != nil {
		t.Error("Error Starting device:", err)
	}
	result, err := ssh.RunSSHCommandWithReply(input)
	if err != nil {
		t.Error("Error executing command:", err)
	}
	err = ssh.Stop()
	if err != nil {
		t.Error("Error Stopping device:", err)
	}
	if result != expected {
		t.Errorf("\"DeviceSsh_Reconnect_WithTimeout Round 1 ('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
	} else {
		t.Logf("\"DeviceSsh_Reconnect_WithTimeout Round 1 ('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
	}

	time.Sleep(15 * time.Second)

	// Second round of connect
	err = ssh.Start()
	if err != nil {
		t.Error("Error Starting device:", err)
	}
	result, err = ssh.RunSSHCommandWithReply(input)
	if err != nil {
		t.Error("Error executing command:", err)
	}
	err = ssh.Stop()
	if err != nil {
		t.Error("Error Stopping device:", err)
	}
	if result != expected {
		t.Errorf("\"DeviceSsh_Reconnect_WithTimeout Round 2 ('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
	} else {
		t.Logf("\"DeviceSsh_Reconnect_WithTimeout Round 2 ('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
	}
}

func Test_DeviceSsh_ReUse_AfterTimeout(t *testing.T) {
	input := "ls /"
	expected := "bin   dev  home  lib32	libx32	mnt  proc  run	 srv  tmp  var\r\nboot  etc  lib	 lib64	media	opt  root  sbin  sys  usr\r\n"
	var err error = nil

	ssh := DeviceSshNew(deviceSshFromYaml())

	// First round of connect
	err = ssh.Start()
	if err != nil {
		t.Error("Error Starting device:", err)
	}
	result, err := ssh.RunSSHCommandWithReply(input)
	if err != nil {
		t.Error("Error executing command:", err)
	}
	if result != expected {
		t.Errorf("\"DeviceSsh_Reconnect_WithTimeout Round 1 ('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
	} else {
		t.Logf("\"DeviceSsh_Reconnect_WithTimeout Round 1 ('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
	}

	time.Sleep(15 * time.Second)

	// Second round of connect
	result, err = ssh.RunSSHCommandWithReply(input)
	if err != nil {
		t.Error("Error executing command:", err)
	}
	err = ssh.Stop()
	if err != nil {
		t.Error("Error Stopping device:", err)
	}
	if result != expected {
		t.Errorf("\"DeviceSsh_Reconnect_WithTimeout Round 2 ('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
	} else {
		t.Logf("\"DeviceSsh_Reconnect_WithTimeout Round 2 ('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
	}
}

func Test_DeviceSsh_ReUse_BeforeTimeout(t *testing.T) {
	input := "ls /"
	expected := "bin   dev  home  lib32	libx32	mnt  proc  run	 srv  tmp  var\r\nboot  etc  lib	 lib64	media	opt  root  sbin  sys  usr\r\n"
	var err error = nil

	ssh := DeviceSshNew(deviceSshFromYaml())

	// First round of connect
	err = ssh.Start()
	if err != nil {
		t.Error("Error Starting device:", err)
	}
	result, err := ssh.RunSSHCommandWithReply(input)
	if err != nil {
		t.Error("Error executing command:", err)
	}
	if result != expected {
		t.Errorf("\"DeviceSsh_Reconnect_WithTimeout Round 1 ('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
	} else {
		t.Logf("\"DeviceSsh_Reconnect_WithTimeout Round 1 ('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
	}

	time.Sleep(5 * time.Second)

	// Second round of connect
	result, err = ssh.RunSSHCommandWithReply(input)
	if err != nil {
		t.Error("Error executing command:", err)
	}
	err = ssh.Stop()
	if err != nil {
		t.Error("Error Stopping device:", err)
	}
	if result != expected {
		t.Errorf("\"DeviceSsh_Reconnect_WithTimeout Round 2 ('%s')\" FAILED, expected -> %v, got -> %v", input, expected, result)
	} else {
		t.Logf("\"DeviceSsh_Reconnect_WithTimeout Round 2 ('%s')\" SUCCEDED, expected -> %v, got -> %v", input, expected, result)
	}
}
