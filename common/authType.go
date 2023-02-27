package common

import (
	"errors"
	"http-knocker/logging"

	"golang.org/x/exp/slices"
)

type AuthType string

const (
	AUTH_TYPE_NONE       AuthType = "none"
	AUTH_TYPE_BASIC_AUTH AuthType = "basic-auth"
	AUTH_TYPE_AUTHELIA   AuthType = "authelia"
)

var (
	authTypeArr = []AuthType{
		AUTH_TYPE_NONE,
		AUTH_TYPE_BASIC_AUTH,
		AUTH_TYPE_AUTHELIA,
	}
)

func (ctx *AuthType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tempStr string = ""
	err := unmarshal(&tempStr)
	if err != nil {
		return err
	}
	tempAuthType := AuthType(tempStr)
	if !slices.Contains(authTypeArr, tempAuthType) {
		logging.CommonLog().Error("Cannot init from string")
		return errors.New("Cannot init from string")
	}
	*ctx = tempAuthType
	return nil
}
