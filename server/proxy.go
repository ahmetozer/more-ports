package server

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)


func proxyHTTP(w http.ResponseWriter, r *http.Request) {
	var remotePort string
	if target, err:= url.Parse("http://"+r.Host); err == nil {
		if tempPort:=target.Port();tempPort == "" {
		remotePort = svConf.defaultPort
		} else {
			remotePort = tempPort
		}
	} else {
		log.Fatal(err)
	}
	origin, err := url.Parse("http://"+svConf.originAddr+":"+remotePort)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Origin addr %s %s",svConf.originAddr,remotePort)
	if err != nil {
		log.Fatal(err)
	}
	p := httputil.NewSingleHostReverseProxy(origin)
	p.ServeHTTP(w, r)
}
