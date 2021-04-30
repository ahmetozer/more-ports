package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	server "github.com/ahmetozer/more-ports/server"
)

var (
	flagVersion bool   = false
	flagHelp    bool   = false
	GitTag      string = ""
	GitCommit   string = ""
	GitUrl      string = ""
	BuildTime   string = ""
	RunningEnv  string = ""
	pidLocation string = ""
	initHelp    string = `
	Sub Commands:
		server:	run more-ports in server mode
`
)

func init() {
	flag.BoolVar(&flagVersion, "v", false, "Show version")
	flag.BoolVar(&flagHelp, "h", false, "Show help")
	flag.StringVar(&pidLocation, "pid", "", "Set pid location")

	flag.Usage = func() {
		flag.PrintDefaults()
		fmt.Println(initHelp)
		pVersion()
	}

	flag.Parse()
	if flagVersion {
		pVersion()
	}
	if flagHelp {
		flag.PrintDefaults()
		fmt.Println(initHelp)
		pVersion()

	}
	if pidLocation != "" {
		generatePidFile()
	}

}

func main() {
	log.Println("More Ports Service")
	args := flag.Args()
	subcmd := ""
	if len(args) > 0 {
		subcmd = args[0]
		args = args[1:]
	} else {
		log.Println("Mode is not given.")
		// subcmd = "server"
		// log.Println("Mode is set to server automaticaly")
	}

	switch subcmd {
	case "server":
		log.Println("Server mode")
		server.RunningEnv = RunningEnv
		server.Main(args)
	case "client":
		log.Println("Client mode")
		//client(args)
	default:
		log.Fatalf("Err mode is not server neither client")
	}
}

func pVersion() {
	fmt.Printf("\n\tMore ports\n")

	if BuildTime != "" {
		fmt.Printf("\tProgram build date: %s\n", BuildTime)
	}

	if GitCommit != "" {
		fmt.Printf("\tCommmit: %s\n", GitCommit)
	}

	if GitTag != "" {
		fmt.Printf("\tTag: %s\n", GitTag)
	}

	if GitUrl != "" {
		fmt.Printf("\tRepo Url: %s\n", GitUrl)
	}

	os.Exit(0)
}
func generatePidFile() {
	pid := []byte(strconv.Itoa(os.Getpid()))
	if err := ioutil.WriteFile(pidLocation, pid, 0644); err != nil {
		log.Fatal(err)
	}
}
