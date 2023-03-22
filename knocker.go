package main

import (
	"errors"
	"os"
	"os/signal"
	"strings"

	"github.com/gorilla/mux"

	"github.com/altucor/http-knocker/common"
	"github.com/altucor/http-knocker/devices"
	"github.com/altucor/http-knocker/firewallCommon"
	"github.com/altucor/http-knocker/firewalls"
	"github.com/altucor/http-knocker/logging"

	auth "github.com/abbot/go-http-auth"
)

type Knocker struct {
	config    common.Configuration
	webServer *WebServer
	devices   map[string]devices.IDevice
	knockers  map[string]*Knock
	router    *mux.Router
	authUsers map[string]string
}

/*
func getEndpointAssociatedWithFirewall(cfg *Configuration, firewallName string) string {
	for key, element := range cfg.Knocks {
		key = key // supres unused key
		if firewallName == element.Firewall {
			return element.Endpoint
		}
	}
	return ""
}
*/

func knockInitDevice(firewallCfg common.FirewallCfg, device devices.IDevice, endpoint common.EndpointCfg) *Knock {
	var knockObject *Knock = nil

	switch firewallCfg.FirewallType.GetValue() {
	case firewallCommon.FIREWALL_BASIC:
		knockObject = KnockNew(
			firewalls.FirewallBasicNew(device, endpoint, firewallCfg),
		)
	}
	return knockObject
}

func parseHtpasswdUserLine(line string) (string, string, error) {
	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return "", "", errors.New("cannot parse htpasswd line")
	}
	return parts[0], parts[1], nil
}

func (ctx *Knocker) setHtpasswdUsersFromArray(users []string) error {
	for _, line := range users {
		user, passHash, err := parseHtpasswdUserLine(line)
		if err != nil {
			logging.CommonLog().Error("cannot parse htpasswd line")
			return errors.New("cannot parse htpasswd line")
		}
		ctx.authUsers[user] = passHash
	}
	return nil
}

func (ctx *Knocker) setHtpasswdUsersFromFile(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		logging.CommonLog().Errorf("cannot parse htpasswd file: %s\n", file)
		return errors.New("cannot parse htpasswd file")
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		user, passHash, err := parseHtpasswdUserLine(line)
		if err != nil {
			logging.CommonLog().Error("cannot parse htpasswd line")
			return errors.New("cannot parse htpasswd line")
		}
		ctx.authUsers[user] = passHash
	}
	return nil
}

func (ctx *Knocker) basicAuthCheck(user string, realm string) string {
	passHash, ok := ctx.authUsers[user]
	if ok {
		return passHash
	}
	return ""
}

func KnockerNew(cfg common.Configuration) *Knocker {
	ctx := &Knocker{
		config:    cfg,
		webServer: NewWebServer(cfg.Server),
		devices:   make(map[string]devices.IDevice),
		knockers:  make(map[string]*Knock),
		router:    mux.NewRouter(),
		authUsers: make(map[string]string),
	}
	for key, dev := range ctx.config.Devices {
		switch dev.Type {
		case devices.DeviceTypeRest:
			ctx.devices[key] = devices.DeviceRestNew(dev.DeviceConnection)
		case devices.DeviceTypeSsh:
			ctx.devices[key] = devices.DeviceSshNew(dev.DeviceConnection)
		case devices.DeviceTypePuller:
			ctx.devices[key] = devices.DevicePullerNew(dev.DeviceConnection)
		}
	}
	for key, element := range ctx.config.Knocks {
		logging.CommonLog().Info("Key:%s => Element:%+v\n", key, element)
		// TODO: Separate devices in another structure
		// To be able initialize devices once and share them
		// between firewalls.
		// Make possible to assign any device to any firewall
		ctx.knockers[key] = knockInitDevice(
			*ctx.config.Firewalls[element.Firewall],
			ctx.devices[ctx.config.Firewalls[element.Firewall].Device],
			*ctx.config.Endpoints[element.Endpoint],
		)
		if ctx.knockers[key] != nil {
			switch ctx.config.Endpoints[element.Endpoint].Auth.Type {
			case common.AUTH_TYPE_NONE:
				ctx.router.HandleFunc(
					"/"+ctx.config.Endpoints[element.Endpoint].Url,
					ctx.knockers[key].GetHttpCallback,
				)
			case common.AUTH_TYPE_BASIC_AUTH:
				if len(ctx.config.Endpoints[element.Endpoint].Auth.Users) != 0 {
					err := ctx.setHtpasswdUsersFromArray(
						ctx.config.Endpoints[element.Endpoint].Auth.Users)
					if err != nil {
						logging.CommonLog().Fatal("Cannot process htpassd users array")
					}
				}
				if ctx.config.Endpoints[element.Endpoint].Auth.UsersFile != "" {
					err := ctx.setHtpasswdUsersFromFile(
						ctx.config.Endpoints[element.Endpoint].Auth.UsersFile)
					if err != nil {
						logging.CommonLog().Fatal("Cannot process htpassd users file")
					}
				}
				if len(ctx.authUsers) == 0 {
					logging.CommonLog().Fatalf("Basic auth users list is empty")
				}
				// TODO: Read more about "realm" and maybe change it to something other
				authenticator := auth.NewBasicAuthenticator("http-knocker", ctx.basicAuthCheck)
				ctx.router.HandleFunc(
					"/"+ctx.config.Endpoints[element.Endpoint].Url,
					authenticator.Wrap(ctx.knockers[key].GetHttpCallbackBasicAuth),
				)
			case common.AUTH_TYPE_AUTHELIA:
				logging.CommonLog().Fatal("Fatal authelia authenticator not implemented")
			default:
				logging.CommonLog().Fatal("Cannot determine AUTH_TYPE for enpoint: %s\n", element.Endpoint)
			}
		}
	}
	return ctx
}

func (ctx *Knocker) Start() {
	for _, item := range ctx.devices {
		item.Start()
	}
	ctx.webServer.Start(ctx.router)
}

func (ctx *Knocker) Wait() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func (ctx *Knocker) Stop() {
	ctx.webServer.Stop()
	for _, item := range ctx.devices {
		item.Stop()
	}
}
