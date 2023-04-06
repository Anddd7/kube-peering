package tunnel

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

type TunnelMode int8

const (
	// client => server, used for vpn
	Forward TunnelMode = iota
	// server => client, used for port-forward (expose)
	Reverse
)

func (t TunnelMode) String() string {
	switch t {
	case Forward:
		return "forward"
	case Reverse:
		return "reverse"
	default:
		return "unknown"
	}
}

func pushTunnelHeaders(req *http.Request, host string) {
	req.Header.Set("X-Forwarded-Host", host)
}

func popTunnelHeaders(req *http.Request) string {
	host := req.Header.Get("X-Forwarded-Host")

	defer func(r *http.Request) {
		r.Header.Del("X-Forwarded-Host")
	}(req)

	return host
}
