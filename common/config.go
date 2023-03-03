package common

import (
	"io/ioutil"

	"github.com/altucor/http-knocker/devices"
	"github.com/altucor/http-knocker/firewallCommon"
	"github.com/altucor/http-knocker/firewallCommon/firewallField"
	"github.com/altucor/http-knocker/logging"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Host                string `yaml:"host"`
	Port                uint16 `yaml:"port"`
	DefaultResponseCode uint16 `yaml:"default-response-code"`
}

type FirewallCfg struct {
	FirewallType    firewallCommon.FirewallType `yaml:"firewallType"`
	DropRuleComment string                      `yaml:"dropRuleCommnet"`
	Protocol        devices.DeviceProtocol      `yaml:"protocol"`
	Device          string                      `yaml:"device"`
}

type IpAddrSource struct {
	Type      IpSourceType `yaml:"type"`
	FieldName string       `yaml:"field-name"`
}

type Auth struct {
	Type      AuthType `yaml:"auth-type"`
	UsersFile string   `yaml:"users-file"`
	Users     []string `yaml:"users"`
}

type EndpointCfg struct {
	IpAddrSource          IpAddrSource           `yaml:"ip-source"`
	Auth                  Auth                   `yaml:"auth"`
	Url                   string                 `yaml:"url"`
	DurationSeconds       DurationSeconds        `yaml:"duration"`
	Port                  uint16                 `yaml:"port"`
	Protocol              firewallField.Protocol `yaml:"protocol"`
	ResponseCodeOnSuccess uint16                 `yaml:"response-code-on-success"`
}

type KnockCfg struct {
	Firewall string `yaml:"firewall"`
	Endpoint string `yaml:"enpoint"`
}

type DeviceWrapper struct {
	Type             devices.DeviceType           `yaml:"type"`
	DeviceConnection devices.DeviceConnectionDesc `yaml:"connection"`
}

type Configuration struct {
	Server    ServerConfig
	Devices   map[string]*DeviceWrapper
	Firewalls map[string]*FirewallCfg
	Endpoints map[string]*EndpointCfg
	Knocks    map[string]*KnockCfg
}

func (ctx *Configuration) SetDefaults() {
	for _, val := range ctx.Endpoints {
		if val.Auth.Type == "" {
			val.Auth.Type = AUTH_TYPE_NONE
		}
		if val.IpAddrSource.Type == "" {
			val.IpAddrSource.Type = IP_SOURCE_TYPE_WEB_SERVER
		}
	}
}

func (ctx *Configuration) Validate() bool {
	for _, val := range ctx.Firewalls {
		if val.Protocol == "" {
			return false
		}
	}
	return true
}

func ConfigurationNew(path string) (Configuration, error) {
	cfg := Configuration{}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		logging.CommonLog().Error("[Config] Error reading file: %s\n", err)
		return cfg, err
	}
	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		logging.CommonLog().Error("[Config] Error unmarshaling yaml file: %s\n", err)
		return cfg, err
	}
	cfg.SetDefaults()
	if !cfg.Validate() {
		logging.CommonLog().Fatal("[Config] Validation of yaml file is failed")
	}
	return cfg, nil
}
