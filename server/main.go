package server

import (
	"flag"
	"log"
	"net/http"

	shared "github.com/ahmetozer/more-ports/pkg"
)

type ServerConfig struct {
	listen      string
	defaultPort string
	remoteAddr  string
	serverName  string
	clientCert  string
	serverCert  string
	serverKey   string
}

var (
	RunningEnv string = ""
	svConf     ServerConfig
	argHelp    bool = false
)

func Main(args []string) {
	//flags := flag.NewFlagSet("server", flag.ContinueOnError)
	svConf.remoteAddr = "127.0.0.1"
	var err error
	if RunningEnv == "container" {
		a, err := listUpInterfaces()
		if err != nil {
			log.Fatal("Err while listing network interfaces %v", err)
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
	flags.StringVar(&svConf.listen, "listen", ":443", "Server listen port")
	flags.StringVar(&svConf.defaultPort, "default-port", "8080", "Origin forward port")
	flags.StringVar(&svConf.remoteAddr, "remote", svConf.remoteAddr, "Remote address for forwarded ports")
	flags.StringVar(&svConf.serverName, "server-name", "", "Server name check from TLS")
	flags.StringVar(&svConf.clientCert, "client-cert", "", "Allow only given client certificate")
	flags.StringVar(&svConf.serverCert, "server-cert", "", "Server cert")
	flags.StringVar(&svConf.serverKey, "server-key", "", "Server key")
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
	certConfig := shared.CertConfig{
		KeyLocation:  currcertDir + "/key.pem",
		CertLocation: currcertDir + "/cert.pem",
	}

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

	httpServer := &http.Server{
		Addr: svConf.listen,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			svConf.proxyHTTP(w, r)
		}),
		// Disable HTTP/2.
		//TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	if svConf.serverName == "" {
		log.Printf("WARN: Flag server-name is not set. The system allows all server names.")
	}
	log.Printf("Starting Server HTTPS server %s\n", svConf.listen)
	log.Fatal(httpServer.ListenAndServeTLS(certConfig.CertLocation, certConfig.KeyLocation))
}
