package devices

import (
	"errors"

	"github.com/altucor/http-knocker/logging"

	"golang.org/x/exp/slices"
)

type DeviceProtocol string

const (
	PROTOCOL_UNKNOWN        DeviceProtocol = "<UNKNOWN>"
	PROTOCOL_ANY            DeviceProtocol = "any"
	PROTOCOL_IP_TABLES      DeviceProtocol = "iptables"
	PROTOCOL_ROUTER_OS_REST DeviceProtocol = "router-os-rest"
	PROTOCOL_ROUTER_OS_SHH  DeviceProtocol = "router-os-ssh"
	PROTOCOL_ROUTER_OS_API  DeviceProtocol = "router-os-api"
)

var (
	protocolArr = []DeviceProtocol{
		PROTOCOL_UNKNOWN,
		PROTOCOL_ANY,
		PROTOCOL_IP_TABLES,
		PROTOCOL_ROUTER_OS_REST,
		PROTOCOL_ROUTER_OS_SHH,
		PROTOCOL_ROUTER_OS_API,
	}
)

func (ctx *DeviceProtocol) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tempStr string = ""
	err := unmarshal(&tempStr)
	if err != nil {
		return err
	}
	tempProtocol := DeviceProtocol(tempStr)
	if !slices.Contains(protocolArr, tempProtocol) {
		logging.CommonLog().Error("Cannot init from string")
		return errors.New("Cannot init from string")
	}
	*ctx = tempProtocol
	return nil
}
