package pkg

import (
	"fmt"
	"log"
	"os"
)

func CertCheck() (string, error) {

	var err error
	certCheckOk := false
	certDir := "/tmp/cert"
	if _, err := os.Stat("/config/key.pem"); err == nil {
		log.Printf("/config/key.pem exists\n")
		certCheckOk = true
	} else {
		log.Printf("/config/key.pem not exist\n")
		certCheckOk = false
	}
	if certCheckOk {
		if _, err := os.Stat("/config/cert.pem"); err == nil {
			fmt.Printf("/config/cert.pem exists\n")
			certCheckOk = true
		} else {
			fmt.Printf("/config/cert.pem not exist\n")
			certCheckOk = false
		}
	}

	if certCheckOk {
		certDir = "/config"
	} else {
		log.Printf("Self created certs will be used\n")
		certDir = "/tmp/cert"
		systemCert := CertConfig{}
		systemCert.Defaults()
		err = systemCert.Generate()
		if err != nil {
			return "", err
		}
	}
	return certDir, err

}
