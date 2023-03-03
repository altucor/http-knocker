package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/altucor/http-knocker/common"
	"github.com/altucor/http-knocker/logging"

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
	s.server.Handler = router
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			logging.CommonLog().Error(err)
		}
	}()
}

func (s *WebServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	s.server.Shutdown(ctx)
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
