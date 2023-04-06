package endpoint

import (
	"errors"
	"net/http"
	"strings"

	"github.com/altucor/http-knocker/firewallCommon/firewallField"
	"github.com/altucor/http-knocker/logging"

	"golang.org/x/exp/slices"
)

type IpSourceType string

const (
	IP_SOURCE_TYPE_WEB_SERVER         IpSourceType = "web-server"
	IP_SOURCE_TYPE_HTTP_HEADERS       IpSourceType = "http-headers"
	IP_SOURCE_TYPE_HTTP_REQUEST_PARAM IpSourceType = "http-request-param"
)

var (
	ipSourceTypeArr = []IpSourceType{
		IP_SOURCE_TYPE_WEB_SERVER,
		IP_SOURCE_TYPE_HTTP_HEADERS,
		IP_SOURCE_TYPE_HTTP_REQUEST_PARAM,
	}
)

func (ctx *IpSourceType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tempStr string = ""
	err := unmarshal(&tempStr)
	if err != nil {
		return err
	}
	tempSourceType := IpSourceType(tempStr)
	if !slices.Contains(ipSourceTypeArr, tempSourceType) {
		logging.CommonLog().Error("Cannot init from string")
		return errors.New("cannot init from string")
	}
	*ctx = tempSourceType
	return nil
}

type IpAddrSource struct {
	Type      IpSourceType `yaml:"type"`
	FieldName string       `yaml:"field-name"`
}

func (ctx *IpAddrSource) SetDefaults() {
	if ctx.Type == "" {
		ctx.Type = IP_SOURCE_TYPE_WEB_SERVER
	}
}

func (ctx *IpAddrSource) GetFromRequest(r *http.Request) (firewallField.Address, error) {
	var clientAddrStr string = ""
	switch ctx.Type {
	case IP_SOURCE_TYPE_WEB_SERVER:
		clientAddrStr = strings.Split(r.RemoteAddr, ":")[0]
	case IP_SOURCE_TYPE_HTTP_HEADERS:
		clientAddrStr = r.Header.Get(ctx.FieldName)
	case IP_SOURCE_TYPE_HTTP_REQUEST_PARAM:
		clientAddrStr = r.URL.Query().Get(ctx.FieldName)
	}

	logging.CommonLog().Debugf("Client addr str: %s", clientAddrStr)
	return firewallField.AddressTypeFromString(clientAddrStr)
}
