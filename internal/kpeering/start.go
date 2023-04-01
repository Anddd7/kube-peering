package kpeering

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/kube-peering/internal/pkg/connectors"
	"github.com/kube-peering/internal/pkg/logger"
	"github.com/kube-peering/internal/pkg/model"
)

type Kpeering struct {
	Interceptor model.Interceptor
	Tunnel      model.Tunnel
}

func (cfg *Kpeering) Start() {
	if cfg.Interceptor.IsTCP() {
		cfg.startTCP()
	}
	if cfg.Interceptor.IsHTTP() {
		cfg.startHttp()
	}
}
func (cfg *Kpeering) startTCP() {
	ctx := context.Background()
	reqChan := make(chan []byte)
	resChan := make(chan []byte)

	interceptor := connectors.NewTCPInterceptor(ctx, cfg.Interceptor, reqChan, resChan)
	tunnel := connectors.NewTunnelServer(ctx, cfg.Tunnel, reqChan, resChan)

	go interceptor.Run()
	go tunnel.Run()

	<-ctx.Done()
}

func (cfg *Kpeering) startHttp() {
	ctx := context.Background()

	tunnel := connectors.NewHttp2TunnelServer(ctx, cfg.Tunnel)
	interceptor := connectors.NewHttp2Interceptor(ctx, cfg.Interceptor)
	handler := func(w http.ResponseWriter, r *http.Request) {
		address := cfg.Tunnel.Address()
		port := address[len(address)-4:]
		url := "https://localhost:" + fmt.Sprint(port)

		req, err := http.NewRequest(http.MethodGet, url, r.Body)
		if err != nil {
			logger.Z.Errorf("error creating request: %v", err)
			return
		}
		resp, err := tunnel.ClientConn.RoundTrip(req)
		if err != nil {
			logger.Z.Errorf("error sending request: %v", err)
			return
		}

		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}

	go tunnel.Run()
	go interceptor.Run(handler)

	<-ctx.Done()
}
