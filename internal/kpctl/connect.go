package kpctl

import (
	"github.com/kube-peering/internal/pkg"
)

type Kpctl struct {
	Tunnel    pkg.Tunnel
	Forwarder pkg.Forwarder
}

func (ctl *Kpctl) Connect() {
	// TODO
}
