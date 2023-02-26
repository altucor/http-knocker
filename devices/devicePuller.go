package devices

import (
	"context"
	"errors"
	"fmt"
	"httpKnocker/deviceCommon"
	"httpKnocker/logging"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

type ConnectionPuller struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Port     uint16 `yaml:"port"`
	Enpoint  string `yaml:"endpoint"`
}

type DevicePuller struct {
	config             ConnectionPuller
	supportedProtocols []DeviceProtocol
	server             *http.Server
	router             *mux.Router
}

func DevicePullerNew(cfg ConnectionPuller) *DevicePuller {
	ctx := &DevicePuller{
		config: cfg,
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

	var wait time.Duration
	ctx.server.Handler = ctx.router
	go func() {
		if err := ctx.server.ListenAndServe(); err != nil {
			logging.CommonLog().Error(err)
		}
	}()
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctxTimeout, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	ctx.server.Shutdown(ctxTimeout)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logging.CommonLog().Info("shutting down")
	os.Exit(0)
	return nil
}

func (ctx *DevicePuller) Stop() error {
	return nil
}

func (ctx *DevicePuller) RunCommandWithReply(command deviceCommon.IDeviceCommand, proto DeviceProtocol) (deviceCommon.IDeviceResponse, error) {
	logging.CommonLog().Error("Not implemented")
	return nil, errors.New("Not Implemented")
}
