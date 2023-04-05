package main

import (
	"os"

	"github.com/kube-peering/internal/pkg"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "http" {
		http()
	} else {
		tcp()
	}
}

func tcp() {
	interceptor := pkg.NewInterceptor("tcp", 10021)
	forwarder := pkg.NewFowarder("tcp", ":8080")
	interceptor.OnTCPConnected = forwarder.ForwardTCP

	interceptor.Start()
}

func http() {
	interceptor := pkg.NewInterceptor("http", 10021)
	forwarder := pkg.NewFowarder("http", "http://localhost:8080")
	interceptor.OnHTTPConnected = forwarder.ForwardHttp
	interceptor.Start()
}
