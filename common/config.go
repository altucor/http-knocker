package common

import (
	"io/ioutil"

	"http-knocker/devices"
	"http-knocker/firewallCommon"
	"http-knocker/firewallCommon/firewallField"
	"http-knocker/logging"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Host                string `yaml:"host"`
	Port                uint16 `yaml:"port"`
	DefaultResponseCode uint16 `yaml:"default-response-code"`
}

type FirewallCfg struct {
	FirewallType      firewallCommon.FirewallType    `yaml:"firewallType"`
	DropRuleComment   string                         `yaml:"dropRuleCommnet"`
	Protocol          devices.DeviceProtocol         `yaml:"protocol"`
	DeviceRouterOsApi *devices.ConnectionRouterOsApi `yaml:"deviceRouterOsApi"`
	DeviceRest        *devices.ConnectionRest        `yaml:"deviceRest"`
	DeviceSsh         *devices.ConnectionSSHCfg      `yaml:"deviceSsh"`
	DevicePuller      *devices.ConnectionPuller      `yaml:"devicePuller"`
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

type Configuration struct {
	Server    ServerConfig
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

func (cfg *Configuration) ConfigDebugPrint() {
	logging.CommonLog().Debug("%+v\n", cfg)
	for key, element := range cfg.Firewalls {
		logging.CommonLog().Debugf("Firewall Key:%s => Element:%+v\n", key, element)
		if element.DeviceRest != nil {
			logging.CommonLog().Debugf("RouterOsRest = %+v\n", element.DeviceRest)
		} else if element.DeviceSsh != nil {
			logging.CommonLog().Debugf("Ssh = %+v\n", element.DeviceSsh)
		}
	}
	for key, element := range cfg.Endpoints {
		logging.CommonLog().Debugf("Endpoint Key:%s => Element:%+v\n", key, element)
	}
	for key, element := range cfg.Knocks {
		logging.CommonLog().Debugf("Knock Key:%s => Element:%+v\n", key, element)
	}
}
