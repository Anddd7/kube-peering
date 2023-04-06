package main

import (
	"os"

	"github.com/kube-peering/internal/pkg"
	"github.com/kube-peering/internal/pkg/connectors"

	example "github.com/kube-peering/example"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "http" {
		// not working, http interceptor cannot accept https request
		connectors.NewProxy(pkg.HTTP, example.ProxyPort, example.AppHttpsAddr).Start()
	} else {
		connectors.NewProxy(pkg.TCP, example.ProxyPort, example.AppHttpsAddr).Start()
	}
}
