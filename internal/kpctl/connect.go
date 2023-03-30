package kpctl

import (
	"context"

	"github.com/kube-peering/internal/pkg/io"
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

	tunnel := io.NewTunnelClient(ctx, ctl.Tunnel, reqChan, resChan)
	forwarder := io.NewTCPForwarder(ctx, ctl.Forwarder, reqChan, resChan)

	go tunnel.Run()
	go forwarder.Run()

	<-ctx.Done()
}
