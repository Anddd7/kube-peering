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
	forwarder := pkg.NewForwarder("tcp", ":8443")
	interceptor.OnTCPConnected = forwarder.ForwardTCP

	interceptor.Start()
}

// not working, http interceptor cannot accept https request
func http() {
	interceptor := pkg.NewInterceptor("http", 10021)
	forwarder := pkg.NewForwarder("http", "https://localhost:8443")
	interceptor.OnHTTPConnected = forwarder.ForwardHTTP
	interceptor.Start()
}
