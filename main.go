package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

var (
	flagServer  bool   = false
	flagClient  bool   = false
	flagVersion bool   = false
	flagHelp    bool   = false
	GitTag      string = ""
	GitCommit   string = ""
	GitUrl      string = ""
	BuildTime   string = ""
	RunningEnv  string = ""
	pidLocation string = ""
)

func init() {
	flag.BoolVar(&flagClient, "client", false, "Run as client mode")
	flag.BoolVar(&flagServer, "server", false, "Run as server mode")
	flag.BoolVar(&flagVersion, "v", false, "Show version")
	flag.BoolVar(&flagHelp, "h", false, "Show help")
	flag.StringVar(&pidLocation, "pid", "/var/run/more-ports.pid", "Set pid location")

	flag.Parse()
	if flagVersion {
		pVersion()
	}
	if flagHelp {
		flag.PrintDefaults()
		pVersion()
	}
	generatePidFile()
}

func main() {
	args := flag.Args()
	fmt.Printf("%s", args)
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
