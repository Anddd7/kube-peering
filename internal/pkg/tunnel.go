package pkg

import (
	"crypto/tls"
	"net"
	"net/http"
)

type Tunnel interface {
	Start()

	SetOnTCPTunnelIn(func(conn *tls.Conn))
	TunnelTCPOut(from *net.TCPConn)

	SetOnHTTPTunnelIn(http.HandlerFunc)
	TunnelHTTPOut(w http.ResponseWriter, r *http.Request)
}
