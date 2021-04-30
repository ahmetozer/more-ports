package pkg

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	case ed25519.PrivateKey:
		return k.Public().(ed25519.PublicKey)
	default:
		return nil
	}
}

type CertConfig struct {
	ed25519Key bool
	isCA       bool
	rsaBits    int
	ecdsaCurve string
	validFrom  string
	validFor   time.Duration
	Hosts      []string
}

func (config *CertConfig) Defaults() {
	if len(config.Hosts) == 0 {
		hostname, err := os.Hostname()
		if err != nil {
			panic(err)
		}
		config.Hosts = []string{hostname}
	}

	if (config.validFor / 1) == 0 {
		config.validFor = 365 * 24 * time.Hour
	}
	if config.rsaBits == 0 {
		config.rsaBits = 2048
	}

}

func (config CertConfig) Generate() error {
	var err error
	// create Certiface dir
	if _, err = os.Stat("/tmp/cert/"); os.IsNotExist(err) {
		if err = os.MkdirAll("/tmp/cert/", os.ModePerm); err != nil {
			return err
		}
	}

	if _, err := os.Stat("/tmp/cert/key.pem"); err == nil {
		if _, err := os.Stat("/tmp/cert/cert.pem"); err != nil {
			return err
		}
	}

	var priv interface{}

	switch config.ecdsaCurve {
	case "":
		if config.ed25519Key {
			_, priv, err = ed25519.GenerateKey(rand.Reader)
			if err != nil {
				return err
			}
		} else {
			priv, err = rsa.GenerateKey(rand.Reader, config.rsaBits)
			if err != nil {
				return err
			}
		}
	case "P224":
		priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "P256":
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "P384":
		priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "P521":
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		return fmt.Errorf("unrecognized elliptic curve: %q", config.ecdsaCurve)
	}
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	var notBefore time.Time
	if len(config.validFrom) == 0 {
		notBefore = time.Now()
	} else {
		notBefore, err = time.Parse("Jan 2 15:04:05 2006", config.validFrom)
		if err != nil {
			return fmt.Errorf("failed to parse creation date: %v", err)
		}
	}

	notAfter := notBefore.Add(config.validFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	for _, h := range config.Hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	if config.isCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %v", err)
	}

	certOut, err := os.Create("/tmp/cert/cert.pem")
	if err != nil {
		return fmt.Errorf("failed to open cert.pem for writing: %v", err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return fmt.Errorf("failed to write data to cert.pem: %v", err)
	}
	if err := certOut.Close(); err != nil {
		return fmt.Errorf("error closing cert.pem: %v", err)
	}

	keyOut, err := os.OpenFile("/tmp/cert/key.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open key.pem for writing: %v", err)
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return fmt.Errorf("unable to marshal private key: %v", err)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		return fmt.Errorf("failed to write data to key.pem: %v", err)
	}
	if err := keyOut.Close(); err != nil {
		return fmt.Errorf("error closing key.pem: %v", err)
	}
	return nil
}
