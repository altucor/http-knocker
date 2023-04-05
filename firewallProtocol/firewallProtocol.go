package firewallProtocol

import (
	"errors"

	"github.com/altucor/http-knocker/logging"

	"golang.org/x/exp/slices"
)

type FirewallProtocolName string

const (
	PROTOCOL_UNKNOWN        FirewallProtocolName = "<UNKNOWN>"
	PROTOCOL_ANY            FirewallProtocolName = "any"
	PROTOCOL_IP_TABLES      FirewallProtocolName = "iptables"
	PROTOCOL_ROUTER_OS_REST FirewallProtocolName = "router-os-rest"
	PROTOCOL_ROUTER_OS_SHH  FirewallProtocolName = "router-os-ssh"
	PROTOCOL_ROUTER_OS_API  FirewallProtocolName = "router-os-api"
)

var (
	protocolArr = []FirewallProtocolName{
		PROTOCOL_UNKNOWN,
		PROTOCOL_ANY,
		PROTOCOL_IP_TABLES,
		PROTOCOL_ROUTER_OS_REST,
		PROTOCOL_ROUTER_OS_SHH,
		PROTOCOL_ROUTER_OS_API,
	}
)

func (ctx *FirewallProtocolName) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tempStr string = ""
	err := unmarshal(&tempStr)
	if err != nil {
		return err
	}
	tempProtocol := FirewallProtocolName(tempStr)
	if !slices.Contains(protocolArr, tempProtocol) {
		logging.CommonLog().Error("Cannot init from string")
		return errors.New("Cannot init from string")
	}
	*ctx = tempProtocol
	return nil
}

type IFirewallProtocol interface {
	GetType() FirewallProtocolName
}
