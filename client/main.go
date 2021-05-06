package client

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
)

var (
	argHelp    bool   = false
	RunningEnv string = ""
	listen     string = ""
	httpPort   string = ""
	httpsPort  string = ""
)

func Main(args []string) {

	flags := flag.NewFlagSet("client", flag.ExitOnError)
	flags.StringVar(&listen, "listen", "localhost:8080", "Proxy server listen addr")
	flags.StringVar(&httpPort, "http-port", ":80", "Remote server port for http requests")
	flags.StringVar(&httpsPort, "https-port", ":443", "Remote server port for https requests")
	flags.BoolVar(&argHelp, "h", false, "Print help for server mode")

	flags.Parse(args)
	if argHelp {
		flags.PrintDefaults()
		return
	}
	log.Printf("Remote ports for http %s, https %s\n", httpPort, httpsPort)
	server := &http.Server{
		ConnContext: SaveConnInContext,
		Addr:        listen,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect { //HTTPS
				handleHTTPConnect(w, r)
			} else { //HTTP
				handleHTTP(w, r)
			}
		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	log.Printf("Client proxy server started at %s\n", listen)
	log.Fatal(server.ListenAndServe())
}
