package main

import (
	"os"

	"github.com/kube-peering/internal/pkg/transit"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "http" {
		http()
	} else {
		tcp()
	}
}

func tcp() {
	interceptor := transit.NewInterceptor("tcp", 10021)
	forwarder := transit.NewFowarder("tcp", ":8443")
	interceptor.OnTCPConnected = forwarder.ForwardTCP

	interceptor.Start()
}

// not working, http interceptor cannot accept https request
func http() {
	interceptor := transit.NewInterceptor("http", 10021)
	forwarder := transit.NewFowarder("http", "https://localhost:8443")
	interceptor.OnHTTPConnected = forwarder.ForwardHttp
	interceptor.Start()
}
