package server

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	shared "github.com/ahmetozer/more-ports/pkg"
)

type ServerConfig struct {
	httpsAddr   string
	httpAddr    string
	defaultPort string
	remoteAddr  string
	serverName  string
	clientCert  string
	serverCert  string
	serverKey   string
}

var (
	argHelp       bool   = false
	httpsRedirect bool   = false
	RunningEnv    string = ""
	svConf        ServerConfig
)

func Main(args []string) {
	//flags := flag.NewFlagSet("server", flag.ContinueOnError)
	svConf.remoteAddr = "127.0.0.1"
	var err error
	if RunningEnv == "container" {
		a, err := listUpInterfaces()
		if err != nil {
			log.Fatalf("Err while listing network interfaces %s", err)
		}
		//Detect network host mode
		if len := len(a); len > 2 {
			log.Printf("You are in a container %v\n", len)
		} else {
			log.Printf("Detecting Origin\n")
			if svConf.remoteAddr, err = DefaultRoute(); err != nil {
				log.Fatalf("%s", err)
			}
			log.Printf("Origin %s\n", svConf.remoteAddr)
		}
	}

	flags := flag.NewFlagSet("server", flag.ExitOnError)
	flags.StringVar(&svConf.httpsAddr, "https", ":443", "HTTPS server listen address")
	flags.StringVar(&svConf.httpAddr, "http", ":80", "HTTP server listen address")
	flags.StringVar(&svConf.defaultPort, "default-port", "8080", "Origin forward port")
	flags.StringVar(&svConf.remoteAddr, "remote", svConf.remoteAddr, "Remote address for forwarded ports")
	flags.StringVar(&svConf.serverName, "server-name", "", "Server name check from TLS")
	flags.StringVar(&svConf.clientCert, "client-cert", "", "Allow only given client certificate")
	flags.StringVar(&svConf.serverCert, "server-cert", "", "Server cert")
	flags.StringVar(&svConf.serverKey, "server-key", "", "Server key")
	flags.BoolVar(&httpsRedirect, "https-redirect", false, "Redirect all http request to https")
	flags.BoolVar(&argHelp, "h", false, "Print help for server mode")
	flags.Parse(args)
	if argHelp {
		flags.PrintDefaults()
		return
	}
	var currcertDir string

	if RunningEnv != "container" {
		currcertDir = "."
	} else {
		currcertDir = "/tmp/cert"
	}

	//svConf.clientCert.Hosts
	certConfig := shared.CertConfig{}

	if svConf.serverCert == "" {
		certConfig.CertLocation = currcertDir + "/cert.pem"
	}
	if svConf.serverKey == "" {
		certConfig.KeyLocation = currcertDir + "/key.pem"
	}

	if svConf.serverName != "" {
		certConfig.Hosts = append(certConfig.Hosts, svConf.serverName)
	}

	err = certConfig.CertCheck()
	if err != nil {
		log.Fatalf("Err while creating Cert %s", err)
	}

	httpsServer := http.Server{
		Addr: svConf.httpsAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			svConf.HTTPStoHTTP(w, r)
		}),
		// Disable HTTP/2.
		//TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	httpServer := http.Server{
		Addr: svConf.httpAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			svConf.HTTPtoHTTP(w, r)
		}),
	}

	if httpsRedirect {
		httpServer.Handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			target := "https://" + req.Host + req.URL.Path
			if len(req.URL.RawQuery) > 0 {
				target += "?" + req.URL.RawQuery
			}
			http.Redirect(w, req, target, http.StatusTemporaryRedirect)
		})
	}

	if svConf.clientCert != "" {
		caCert, err := ioutil.ReadFile(svConf.clientCert)
		if err != nil {
			log.Fatal(err)
		}

		roots := x509.NewCertPool()
		ok := roots.AppendCertsFromPEM(caCert)
		if !ok {
			log.Fatal("failed to parse root certificate")
		}
		httpsServer.TLSConfig = &tls.Config{
			ClientAuth: tls.RequireAndVerifyClientCert,
			ClientCAs:  roots,
		}

	}

	if svConf.serverName == "" {
		log.Printf("WARN: Flag server-name is not set. The system allows all server names.")
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		log.Printf("Starting Server HTTPS server %s\n", svConf.httpsAddr)
		log.Fatal(httpsServer.ListenAndServeTLS(certConfig.CertLocation, certConfig.KeyLocation))
		wg.Done()
	}()

	go func() {
		log.Printf("Starting Server HTTP server %s\n", svConf.httpAddr)
		log.Fatal(httpServer.ListenAndServe())
		wg.Done()
	}()

	wg.Wait()

}
