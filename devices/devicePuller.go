package devices

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/altucor/http-knocker/deviceCommon"
	"github.com/altucor/http-knocker/logging"

	"github.com/gorilla/mux"
)

type ConnectionPuller struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Port     uint16 `yaml:"port"`
	Endpoint string `yaml:"endpoint"`
}

type DevicePuller struct {
	config             ConnectionPuller
	supportedProtocols []DeviceProtocol
	server             *http.Server
	router             *mux.Router
}

func (ctx *DevicePuller) defaultHandler(w http.ResponseWriter, r *http.Request) {
	logging.CommonLog().Info("[devicePuller] called defaultHandler")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404\n")
}

func DevicePullerNew(cfg DeviceConnectionDesc) *DevicePuller {
	ctx := &DevicePuller{
		config: ConnectionPuller{
			Username: cfg.Username,
			Password: cfg.Password,
			Port:     cfg.Port,
			Endpoint: cfg.Endpoint,
		},
		supportedProtocols: []DeviceProtocol{
			PROTOCOL_ANY,
		},
		server: &http.Server{
			Addr:         "0.0.0.0" + ":" + fmt.Sprint(cfg.Port),
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		router: mux.NewRouter(),
	}
	ctx.router.HandleFunc(ctx.config.Endpoint, ctx.defaultHandler)
	return ctx
}

func (ctx *DevicePuller) GetSupportedProtocols() []DeviceProtocol {
	return ctx.supportedProtocols
}

func (ctx *DevicePuller) GetType() DeviceType {
	return DeviceTypePuller
}

func http_not_found_handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "1337\n")
}

func (ctx *DevicePuller) Start() error {
	ctx.router.HandleFunc(
		"/",
		http_not_found_handler,
	)
	ctx.server.Handler = ctx.router
	go func() {
		if err := ctx.server.ListenAndServe(); err != nil {
			logging.CommonLog().Error(err)
		}
	}()
	return nil
}

func (ctx *DevicePuller) Stop() error {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	ctx.server.Shutdown(ctxTimeout)
	return nil
}

func (ctx *DevicePuller) RunCommandWithReply(command deviceCommon.IDeviceCommand, proto DeviceProtocol) (deviceCommon.IDeviceResponse, error) {
	logging.CommonLog().Error("Not implemented")
	return nil, errors.New("not Implemented")
}
