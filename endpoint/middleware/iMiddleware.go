package middleware

import (
	"net/http"

	"gopkg.in/yaml.v3"
)

type iMiddleware interface {
	Register(handler http.HandlerFunc) http.HandlerFunc
	Intercept(w http.ResponseWriter, r *http.Request) error
}

type Config struct {
	Type string `yaml:"type"`
}

type InterfaceWrapper struct {
	Middleware iMiddleware
	Config     Config
}

func (ctx *InterfaceWrapper) UnmarshalYAML(value *yaml.Node) error {
	if err := value.Decode(&ctx.Config); err != nil {
		return err
	}

	var err error = nil
	switch ctx.Config.Type {
	case "basic-auth":
		ctx.Middleware, err = BasicAuthNewFromYaml(value)
	}
	if err != nil {
		return err
	}
	return nil
}
