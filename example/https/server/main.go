package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	example "github.com/kube-peering/example/http2"
)

func main() {
	server := http.Server{
		Addr:      example.Port,
		TLSConfig: tlsConfig(),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("Protocol: %s", r.Proto)))
	})

	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatal(err)
	}
}

func tlsConfig() *tls.Config {
	crt, err := ioutil.ReadFile(example.CertPath)
	if err != nil {
		log.Fatal(err)
	}

	key, err := ioutil.ReadFile(example.KeyPath)
	if err != nil {
		log.Fatal(err)
	}

	cert, err := tls.X509KeyPair(crt, key)
	if err != nil {
		log.Fatal(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   example.ServerName,
	}
}
