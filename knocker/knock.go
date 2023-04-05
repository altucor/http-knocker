package knocker

import (
	"fmt"
	"net/http"

	"github.com/altucor/http-knocker/endpoint"
	"github.com/altucor/http-knocker/firewallControllers"
	"github.com/altucor/http-knocker/logging"
	"gopkg.in/yaml.v3"
)

type KnockCfg struct {
	Controller string `yaml:"controller"`
	Endpoint   string `yaml:"enpoint"`
}

type Knock struct {
	cfg        KnockCfg
	controller firewallControllers.IController
	endpoint   *endpoint.Endpoint
}

func (ctx *Knock) UnmarshalYAML(value *yaml.Node) error {
	var cfg KnockCfg
	if err := value.Decode(&cfg); err != nil {
		return err
	}
	ctx.cfg = cfg
	return nil
}

func (ctx *Knock) GetControllerName() string {
	return ctx.cfg.Controller
}

func (ctx *Knock) SetController(controller firewallControllers.IController) {
	ctx.controller = controller
}

func (ctx *Knock) GetEndpointName() string {
	return ctx.cfg.Endpoint
}

func (ctx *Knock) SetEndpoint(endpoint *endpoint.Endpoint) {
	ctx.endpoint = endpoint
}

func (ctx *Knock) GetEndpoint() endpoint.Endpoint {
	return *ctx.endpoint
}

func (ctx *Knock) GetHttpCallback(w http.ResponseWriter, r *http.Request) {
	logging.CommonLog().Info("[knock] accessing knock endpoint:", ctx.endpoint.Url)

	if clientAddr, err := ctx.endpoint.IpAddrSource.GetFromRequest(r); err != nil {
		logging.CommonLog().Error("[knock] Error parsing client address:", err)
	} else {
		// Perform adding client in another thread
		// To be able response to HTTP client faster
		// And prevent timing attacks
		go ctx.controller.AddClient(clientAddr)
		if ctx.endpoint.ResponseCodeOnSuccess != 0 {
			w.WriteHeader(int(ctx.endpoint.ResponseCodeOnSuccess))
			fmt.Fprintf(w, "%d\n", ctx.endpoint.ResponseCodeOnSuccess)
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "404\n")
		}
		//http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
