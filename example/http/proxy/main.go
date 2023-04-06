package main

import (
	"os"

	"github.com/kube-peering/internal/pkg"

	example "github.com/kube-peering/example"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "http" {
		http()
	} else {
		tcp()
	}
}

func tcp() {
	interceptor := pkg.NewInterceptor("tcp", example.ProxyPort)
	forwarder := pkg.NewForwarder("tcp", example.AppAddr)
	interceptor.OnTCPConnected = forwarder.ForwardTCP

	interceptor.Start()
}

func http() {
	interceptor := pkg.NewInterceptor("http", example.ProxyPort)
	forwarder := pkg.NewForwarder("http", example.AppAddr)
	interceptor.OnHTTPConnected = forwarder.ForwardHTTP

	interceptor.Start()
}
