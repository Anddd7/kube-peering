package main

import (
	example "github.com/kube-peering/example"
	"github.com/kube-peering/internal/pkg"
	"github.com/kube-peering/internal/pkg/connectors"
)

func main() {
	connectors.NewProxy(pkg.TCP, example.ProxyPort, example.AppAddr).Start()
}
