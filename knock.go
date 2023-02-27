package main

import (
	"fmt"
	"net/http"
	"strings"

	"http-knocker/common"
	"http-knocker/firewallCommon/firewallField"
	"http-knocker/firewalls"
	"http-knocker/logging"

	auth "github.com/abbot/go-http-auth"
)

type Knock struct {
	firewall firewalls.IFirewall
}

func KnockNew(firewall firewalls.IFirewall) *Knock {
	ctx := &Knock{
		firewall: firewall,
	}
	return ctx
}

func (ctx *Knock) GetHttpCallback(w http.ResponseWriter, r *http.Request) {
	logging.CommonLog().Info("[knock] accessing knock endpoint:", ctx.firewall.GetEndpoint().Url)

	var clientAddrStr string = ""
	switch ctx.firewall.GetEndpoint().IpAddrSource.Type {
	case common.IP_SOURCE_TYPE_WEB_SERVER:
		clientAddrStr = strings.Split(r.RemoteAddr, ":")[0]
	case common.IP_SOURCE_TYPE_HTTP_HEADERS:
		clientAddrStr = r.Header.Get(ctx.firewall.GetEndpoint().IpAddrSource.FieldName)
	case common.IP_SOURCE_TYPE_HTTP_REQUEST_PARAM:
		clientAddrStr = r.URL.Query().Get(ctx.firewall.GetEndpoint().IpAddrSource.FieldName)
	}

	logging.CommonLog().Debugf("Client addr str: %s\n", clientAddrStr)
	clientAddr, err := firewallField.AddressTypeFromString(clientAddrStr)
	if err != nil {
		logging.CommonLog().Error("[knock] Error parsing client address:", err)
	} else {
		// Perform adding client in another thread
		// To be able response to HTTP client faster
		// And prevent timing attacks
		go ctx.firewall.AddClient(clientAddr)
		if ctx.firewall.GetEndpoint().ResponseCodeOnSuccess != 0 {
			w.WriteHeader(int(ctx.firewall.GetEndpoint().ResponseCodeOnSuccess))
			fmt.Fprintf(w, "%d\n", ctx.firewall.GetEndpoint().ResponseCodeOnSuccess)
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "404\n")
		}
		//http.Error(rw, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (ctx *Knock) GetHttpCallbackBasicAuth(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	logging.CommonLog().Info("Basic Auth handler called")
	ctx.GetHttpCallback(w, &r.Request)
}
