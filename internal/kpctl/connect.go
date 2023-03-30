package kpctl

import (
	"context"

	"github.com/kube-peering/internal/pkg/connectors"
	"github.com/kube-peering/internal/pkg/model"
)

type Kpctl struct {
	Tunnel    model.Tunnel
	Forwarder model.Forwarder
}

func (ctl *Kpctl) Connect() {
	ctx := context.Background()
	reqChan := make(chan []byte)
	resChan := make(chan []byte)

	tunnel := connectors.NewTunnelClient(ctx, ctl.Tunnel, reqChan, resChan)
	forwarder := connectors.NewTCPForwarder(ctx, ctl.Forwarder, reqChan, resChan)

	go tunnel.Run()
	go forwarder.Run()

	<-ctx.Done()
}
