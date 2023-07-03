package webserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/altucor/http-knocker/logging"
	"gopkg.in/yaml.v3"

	"github.com/gorilla/mux"
)

type ServerConfig struct {
	Host                string `yaml:"host"`
	Port                uint16 `yaml:"port"`
	DefaultResponseCode uint16 `yaml:"default-response-code"`
}

func http_not_found_handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404\n")
}

type WebServer struct {
	config ServerConfig
	server *http.Server
	router *mux.Router
}

// func WebServerNew(cfg ServerConfig) *WebServer {
// 	srv := &WebServer{
// 		config: cfg,
// 		server: &http.Server{
// 			Addr:         cfg.Host + ":" + fmt.Sprint(cfg.Port),
// 			WriteTimeout: 15 * time.Second,
// 			ReadTimeout:  15 * time.Second,
// 			IdleTimeout:  60 * time.Second,
// 		},
// 		router: mux.NewRouter(),
// 	}
// 	srv.router.NotFoundHandler = http.HandlerFunc(http_not_found_handler)
// 	return srv
// }

func (ctx *WebServer) UnmarshalYAML(value *yaml.Node) error {
	var cfg ServerConfig
	if err := value.Decode(&cfg); err != nil {
		return err
	}
	ctx.config = cfg
	ctx.server = &http.Server{
		Addr:         cfg.Host + ":" + fmt.Sprint(cfg.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	ctx.router = mux.NewRouter()
	ctx.router.NotFoundHandler = http.HandlerFunc(http_not_found_handler)
	return nil
}

func (s *WebServer) Start() {
	logging.CommonLog().Info("[webServer] Starting...")
	s.server.Handler = s.router
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			logging.CommonLog().Error(err)
		}
	}()
	logging.CommonLog().Info("[webServer] Starting... DONE")
}

func (s *WebServer) Stop() {
	logging.CommonLog().Info("[webServer] Stopping...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	s.server.Shutdown(ctx)
	logging.CommonLog().Info("[webServer] Stopping... DONE")
}

/*
func http_firewall_handler(firewall IFirewall) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		firewall.AddClientToAcceptedList("addr_list", r.RemoteAddr, "1h")
	}
}
*/

func (s *WebServer) AddEndpoint(url string, handler http.HandlerFunc) {
	logging.CommonLog().Info("adding server endpoint: ", url)
	s.router.HandleFunc(url, handler)
}
