package pkg

import (
	"net/http"
)

type Tunnel interface {
	Start()

	SetOnTCPTunnel(func(conn PipeConn))
	TunnelTCP(from PipeConn)

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

type TunnelConfig struct {
	Port           int
	Host           string
	ServerCertPath string `toml:"server_cert_path"`
	ServerKeyPath  string `toml:"server_key_path"`
	CaCertPath     string `toml:"ca_cert_path"`
	ServerName     string `toml:"server_name"`
}
