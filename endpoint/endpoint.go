package endpoint

import (
	"crypto/sha1"
	"fmt"
	"net/http"

	"github.com/altucor/http-knocker/common"
	"github.com/altucor/http-knocker/endpoint/middleware"
	"github.com/altucor/http-knocker/firewallCommon/firewallField"
	"gopkg.in/yaml.v3"
)

type Endpoint struct {
	// IpAddrSource - describe where from to take client IP address
	IpAddrSource IpAddrSource `yaml:"ip-source"`
	// Middlewares additionally do some checks for incoming requests
	// Like password protection or DDoS prevention
	Middlewares           map[string]middleware.InterfaceWrapper `yaml:"middlewares"`
	DurationSeconds       common.DurationSeconds                 `yaml:"duration"`
	Port                  uint16                                 `yaml:"port"`
	Protocol              firewallField.Protocol                 `yaml:"protocol"`
	ResponseCodeOnSuccess uint16                                 `yaml:"response-code-on-success"`
}

func EndpointNewFromConfig(value *yaml.Node) (Endpoint, error) {
	var cfg Endpoint
	if err := value.Decode(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (ctx *Endpoint) SetDefaults() {
	ctx.IpAddrSource.SetDefaults()
	// ctx.Auth.SetDefaults()
}

func (ctx *Endpoint) RegisterMiddlewares(input http.HandlerFunc) http.HandlerFunc {
	var last http.HandlerFunc = input
	for _, element := range ctx.Middlewares {
		last = element.Middleware.Register(last)
	}
	return last
}

// func (ctx Endpoint) RegisterWithMiddlewares(final http.HandlerFunc) (string, http.HandlerFunc) {
// 	switch ctx.Auth.Type {
// 	case AUTH_TYPE_NONE:
// 		return "/" + ctx.Url, final
// 	case AUTH_TYPE_BASIC_AUTH:
// 		return "/" + ctx.Url, ctx.Auth.GetAuthenticator(final)
// 	default:
// 		logging.CommonLog().Error("Invalid auth type")
// 		return "/" + ctx.Url, nil
// 	}
// }

func (ctx Endpoint) GetHash(url string) string {
	// Do not use here duration field to allow change timeout for already added rules
	h := sha1.New()
	h.Write([]byte(fmt.Sprint(ctx.Port) + ctx.Protocol.GetString() + url))
	return fmt.Sprintf("%x", h.Sum(nil))
}
