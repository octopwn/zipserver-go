// main.go
package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"os"

	"github.com/octopwn/zipserver-go/server"
)

func main() {
	var address string
	var port int
	var sslCert string
	var sslKey string
	var sslCA string

	flag.StringVar(&address, "a", "0.0.0.0", "IP/hostname to listen on")
	flag.IntVar(&port, "p", 8000, "Port to listen on")
	flag.StringVar(&sslCert, "ssl-cert", "", "Certificate file for SSL")
	flag.StringVar(&sslKey, "ssl-key", "", "Key file for SSL")
	flag.StringVar(&sslCA, "ssl-ca", "", "CA cert file for client cert validations")
	flag.Parse()

	zipFilePath := flag.Arg(0)
	if zipFilePath == "" {
		log.Fatal("Path to the zip file to serve is required")
	}

	var tlsConfig *tls.Config
	if sslCert != "" {
		if sslKey == "" {
			log.Fatal("TLS certificate is set but no keyfile!")
		}

		tlsConfig = &tls.Config{}
		cert, err := tls.LoadX509KeyPair(sslCert, sslKey)
		if err != nil {
			log.Fatalf("Failed to load certificate: %v", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}

		if sslCA != "" {
			caCertPool := x509.NewCertPool()
			caCert, err := os.ReadFile(sslCA)
			if err != nil {
				log.Fatalf("Failed to read CA certificate: %v", err)
			}
			if !caCertPool.AppendCertsFromPEM(caCert) {
				log.Fatal("Failed to append CA certificate")
			}
			tlsConfig.ClientCAs = caCertPool
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}
	}

	server.Serve(zipFilePath, address, port, tlsConfig)
}
