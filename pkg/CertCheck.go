package pkg

import (
	"log"
	"os"
)

func (config *CertConfig) CertCheck() error {

	var err error
	certCheckOk := false

	if _, err := os.Stat(config.CertDir + "/key.pem"); err == nil {
		log.Printf(config.CertDir + "/key.pem exists\n")
		certCheckOk = true
	} else {
		log.Printf(config.CertDir + "/key.pem not exist\n")
		certCheckOk = false
	}
	if certCheckOk {
		if _, err := os.Stat(config.CertDir + "/cert.pem"); err == nil {
			log.Printf(config.CertDir + "/cert.pem exists\n")
			certCheckOk = true
		} else {
			log.Printf(config.CertDir + "cert.pem not exist\n")
			certCheckOk = false
		}
	}

	if !certCheckOk {
		log.Printf("Self created certs will be used\n")
		config.Defaults()
		err = config.Generate()
		if err != nil {
			return err
		}
	}
	return nil

}
