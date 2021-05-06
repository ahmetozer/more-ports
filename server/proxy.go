package server

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func (obj *ServerConfig) HTTPStoHTTP(w http.ResponseWriter, r *http.Request) {
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

func (obj *ServerConfig) HTTPtoHTTP(w http.ResponseWriter, r *http.Request) {
	if obj.serverName != "" {
		if obj.serverName != strings.Split(r.Host, ":")[0] {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "The request server name is invalid. %s", r.Host)
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
