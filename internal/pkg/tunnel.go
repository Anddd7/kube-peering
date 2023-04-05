package pkg

import (
	"net/http"

	"github.com/kube-peering/internal/pkg/util"
)

type Tunnel interface {
	Start()

	SetOnTCPTunnel(func(conn util.PipeConn))
	TunnelTCP(from util.PipeConn)

	SetOnHTTPTunnel(http.HandlerFunc)
	TunnelHTTP(w http.ResponseWriter, r *http.Request)
}
