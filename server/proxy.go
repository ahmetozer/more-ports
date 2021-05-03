package server

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func (obj *ServerConfig) proxyHTTP(w http.ResponseWriter, r *http.Request) {
	if obj.serverName != "" {
		if obj.serverName != r.TLS.ServerName {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "The request server name is invalid.")
			return
		}
	}
	var remotePort string
	if target, err := url.Parse("http://" + r.Host); err == nil {
		if tempPort := target.Port(); tempPort == "" {
			remotePort = svConf.defaultPort
		} else {
			remotePort = tempPort
		}
	} else {
		log.Fatal(err)
	}
	origin, err := url.Parse("http://" + svConf.remoteAddr + ":" + remotePort)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	p := httputil.NewSingleHostReverseProxy(origin)
	p.ServeHTTP(w, r)
}
