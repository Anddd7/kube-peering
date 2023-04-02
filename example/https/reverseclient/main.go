package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/http2"

	example "github.com/kube-peering/example/http2"
)

func main() {
	conn, _ := tls.Dial("tcp", example.Localhost, tlsConfig())
	h2s := &http2.Server{}
	h2s.ServeConn(conn, &http2.ServeConnOpts{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(
				fmt.Sprintf(
					"Reverse connection from client %s => %s, Protocol: %s",
					r.RemoteAddr, conn.LocalAddr(), r.Proto,
				),
			))
		}),
	})
}

func transport2() *http2.Transport {
	return &http2.Transport{
		TLSClientConfig:    tlsConfig(),
		DisableCompression: true,
		AllowHTTP:          false,
	}
}

func tlsConfig() *tls.Config {
	crt, err := ioutil.ReadFile(example.CertPath)
	if err != nil {
		log.Fatal(err)
	}

	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(crt)

	return &tls.Config{
		RootCAs:            rootCAs,
		InsecureSkipVerify: false,
		ServerName:         example.ServerName,
	}
}
