package main

import (
	"fmt"
	"net/http"
	"os"

	example "github.com/kube-peering/example"
	"github.com/kube-peering/internal/pkg"
	"github.com/kube-peering/internal/pkg/tunnel"
)

func main() {
	server()
	client()

	select {}
}

func server() {
	_, port := tunnelPorts()
	tunnel := tunnel.NewTunnelServer("http", 10086, example.TunnelServerCert, example.TunnelServerKey, example.TunnelServerName)

	if port == 0 {
		fowarder := pkg.NewFowarder("http", "localhost:8080")
		tunnel.SetOnHTTPTunnel(func(w http.ResponseWriter, r *http.Request) {
			method := r.Header.Get("X-Forwarded-Method")
			path := r.Header.Get("X-Forwarded-Path")
			url := "http://localhost:8080"

			req, _ := http.NewRequest(
				method,
				fmt.Sprintf("%s%s", url, path),
				r.Body,
			)
			req.Header = r.Header.Clone()

			fmt.Printf("content length: %d", req.ContentLength)

			fowarder.ForwardHTTP(w, req)
		})
	} else {
		// interceptor := pkg.NewInterceptor("tcp", port)
		// interceptor.OnHTTPConnected = tunnel.TunnelHTTPOut
		// go interceptor.Start()
	}

	go tunnel.Start()
}

func client() {
	port, _ := tunnelPorts()
	tunnel := tunnel.NewTunnelClient("http", "localhost:10086", example.TunnelCaCert, example.TunnelServerName)

	if port == 0 {
		// fowarder := pkg.NewFowarder("tcp", ":8080")
		// tunnel.SetOnHTTPTunnelIn(fowarder.ForwardHTTP)
	} else {
		interceptor := pkg.NewInterceptor("http", port)
		interceptor.OnHTTPConnected = tunnel.TunnelHTTP
		go interceptor.Start()
	}

	go tunnel.Start()
}

// client will connect to 10022
// normal : client -> tunnel client --------> tunnel server -> server
// reverse: client -> tunnel server --------> tunnel client -> server
func tunnelPorts() (int, int) {
	if len(os.Args) > 1 {
		if os.Args[1] == "reverse" {
			return 0, 10022
		}
	}
	return 10022, 0
}
