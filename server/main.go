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
	originAddr  string
}

var (
	RunningEnv string = ""
	svConf     ServerConfig
	argHelp    bool = false
)

func Main(args []string) {
	//flags := flag.NewFlagSet("server", flag.ContinueOnError)
	flags := flag.NewFlagSet("server", flag.ExitOnError)
	flags.StringVar(&svConf.listen, "listen", ":443", "Server listen port")
	flags.StringVar(&svConf.defaultPort, "default-port", "8080", "Origin forward port")
	flags.StringVar(&svConf.originAddr, "origin", "127.0.0.1", "Origin address")
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

	certConfig := shared.CertConfig{
		CertDir: currcertDir,
	}

	err := certConfig.CertCheck()
	if err != nil {
		log.Fatalf("Err while creating Cert %s", err)
	}

	httpServer := &http.Server{
		Addr: svConf.listen,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			proxyHTTP(w, r)
		}),
		// Disable HTTP/2.
		//TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	log.Printf("Starting Server HTTPS server\n")
	log.Fatal(httpServer.ListenAndServeTLS(currcertDir+"/cert.pem", currcertDir+"/key.pem"))
}
