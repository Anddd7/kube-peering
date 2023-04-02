package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	example "github.com/kube-peering/example/http2"
	"golang.org/x/net/http2"
)

func main() {
	listener, _ := net.Listen("tcp", example.Port)
	conn, _ := listener.Accept()
	port := conn.LocalAddr().(*net.TCPAddr).Port

	t := &http2.Transport{}
	tlsConn := tls.Server(conn, tlsConfig())
	h2Conn, err := t.NewClientConn(tlsConn)
	if err != nil {
		// tls verification failed
		log.Fatal(err)
	}

	fmt.Println(conn.LocalAddr().String() + " => " + conn.RemoteAddr().String())

	url := "https://localhost:" + fmt.Sprint(port)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := h2Conn.RoundTrip(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	res.Body.Close()

	fmt.Printf("Code: %d\n", res.StatusCode)
	fmt.Printf("Body: %s\n", body)
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
