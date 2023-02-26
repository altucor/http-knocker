package main

import (
	"context"
	"fmt"
	"httpKnocker/common"
	"httpKnocker/logging"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func http_not_found_handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404\n")
}

type WebServer struct {
	config common.ServerConfig
	server *http.Server
}

func NewWebServer(cfg common.ServerConfig) *WebServer {
	srv := &WebServer{
		config: cfg,
		server: &http.Server{
			Addr:         cfg.Host + ":" + fmt.Sprint(cfg.Port),
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
	//srv.router.NotFoundHandler = http.HandlerFunc(http_not_found_handler)
	//srv.router.NotFoundHandler = http.HandlerFunc(http_wrapped_handler(fwcfg))
	//fwcfg.Server.Port = 1337
	return srv
}

func (s *WebServer) Start(router *mux.Router) {
	var wait time.Duration
	s.server.Handler = router
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
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
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	s.server.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logging.CommonLog().Info("shutting down")
	os.Exit(0)
}

/*
func http_firewall_handler(firewall IFirewall) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		firewall.AddClientToAcceptedList("addr_list", r.RemoteAddr, "1h")
	}
}
*/

func (s *WebServer) addEndpoint(url string, handler http.HandlerFunc) {
	logging.CommonLog().Info("adding server endpoint", url)
	//s.router.HandleFunc("/"+url, handler)
	//s.router.Handle("/"+url, )
	//s.router.HandleFunc("/"+endpoint.Url, http_firewall_handler(firewall))
	//s.router.HandleFunc("/"+endpoint.Url, http_listener)
	//s.router.HandleFunc("/"+endpoint.Url+"/", http_listener)
}
