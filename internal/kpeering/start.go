package kpeering

import (
	"github.com/kube-peering/internal/pkg"
)

type Kpeering struct {
	Interceptor pkg.Interceptor
	Tunnel      pkg.Tunnel
}

func (cfg *Kpeering) Start() {
	// TODO
}
