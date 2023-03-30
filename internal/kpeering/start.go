package kpeering

import (
	"context"

	"github.com/kube-peering/internal/pkg/io"
	"github.com/kube-peering/internal/pkg/model"
)

type Kpeering struct {
	Interceptor model.Interceptor
	Tunnel      model.Tunnel
}

func (cfg *Kpeering) Start() {
	ctx := context.Background()
	reqChan := make(chan []byte)
	resChan := make(chan []byte)

	interceptor := io.NewTCPInterceptor(ctx, cfg.Interceptor, reqChan, resChan)
	tunnel := io.NewTunnelServer(ctx, cfg.Tunnel, reqChan, resChan)

	go interceptor.Run()
	go tunnel.Run()

	<-ctx.Done()
}
