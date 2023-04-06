package endpoint

import (
	"crypto/sha1"
	"fmt"
	"net/http"

	"github.com/altucor/http-knocker/common"
	"github.com/altucor/http-knocker/firewallCommon/firewallField"
	"github.com/altucor/http-knocker/logging"
	"gopkg.in/yaml.v3"
)

type Endpoint struct {
	// IpAddrSource - describe where from to take client IP address
	IpAddrSource IpAddrSource `yaml:"ip-source"`
	// Auth - if set then adds client verification mechanism, for example: BasicAuth
	Auth                  Auth                   `yaml:"auth"`
	Url                   string                 `yaml:"url"`
	DurationSeconds       common.DurationSeconds `yaml:"duration"`
	Port                  uint16                 `yaml:"port"`
	Protocol              firewallField.Protocol `yaml:"protocol"`
	ResponseCodeOnSuccess uint16                 `yaml:"response-code-on-success"`
}

func (ctx *Endpoint) SetDefaults() {
	ctx.IpAddrSource.SetDefaults()
	ctx.Auth.SetDefaults()
}

func (ctx *Endpoint) IsEqual(other *Endpoint) bool {
	if ctx.Url == other.Url {
		return true
	}
	// if ctx.Port == other.Port {
	// 	return true
	// }
	// if ctx.Protocol.GetValue() == other.Protocol.GetValue() {
	// 	return true
	// }
	return false
}

func (ctx Endpoint) RegisterWithMiddlewares(final http.HandlerFunc) (string, http.HandlerFunc) {
	switch ctx.Auth.Type {
	case AUTH_TYPE_NONE:
		return "/" + ctx.Url, final
	case AUTH_TYPE_BASIC_AUTH:
		return "/" + ctx.Url, ctx.Auth.GetAuthenticator(final)
	default:
		logging.CommonLog().Error("Invalid auth type")
		return "/" + ctx.Url, nil
	}
}

func (ctx Endpoint) GetHash() string {
	h := sha1.New()
	h.Write([]byte(ctx.Url + fmt.Sprintf("%d", ctx.DurationSeconds) + fmt.Sprint(ctx.Port) + ctx.Protocol.GetString()))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func EndpointNewFromConfig(value *yaml.Node) (Endpoint, error) {
	var cfg Endpoint
	if err := value.Decode(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
