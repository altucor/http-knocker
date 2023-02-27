package common

import (
	"errors"
	"http-knocker/logging"

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
		return errors.New("Cannot init from string")
	}
	*ctx = tempSourceType
	return nil
}
