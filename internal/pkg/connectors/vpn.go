package connectors

import (
	"fmt"

	"github.com/kube-peering/internal/pkg"
	"github.com/kube-peering/internal/pkg/tunnel"
)

type VPNConfig struct {
	Tunnel     TunnelConfig
	Protocol   pkg.Protocol
	LocalPort  int
	RemoteAddr string
}

type TunnelConfig struct {
	Port           int
	Host           string
	ServerCertPath string
	ServerKeyPath  string
	CaCertPath     string
	ServerName     string
}

type VPNServer struct {
	Protocol  pkg.Protocol
	Forwarder *pkg.Forwarder
	Tunnel    tunnel.Tunnel
}

func NewVPNServer(cfg VPNConfig) *VPNServer {
	_tunnel := tunnel.NewTunnelServer(
		tunnel.Forward,
		cfg.Protocol, cfg.Tunnel.Port,
		cfg.Tunnel.ServerCertPath, cfg.Tunnel.ServerKeyPath, cfg.Tunnel.ServerName,
	)
	forwarder := pkg.NewForwarder(cfg.Protocol, cfg.RemoteAddr)

	_tunnel.SetOnTCPTunnel(forwarder.ForwardTCP)
	_tunnel.SetOnHTTPTunnel(forwarder.ForwardHTTP)

	return &VPNServer{
		Protocol:  cfg.Protocol,
		Forwarder: forwarder,
		Tunnel:    _tunnel,
	}
}

func (s *VPNServer) Start() {
	s.Tunnel.Start()
}

type VPNClient struct {
	Protocol    pkg.Protocol
	Interceptor *pkg.Interceptor
	Tunnel      tunnel.Tunnel
}

func NewVPNClient(cfg VPNConfig) *VPNClient {
	interceptor := pkg.NewInterceptor(cfg.Protocol, cfg.LocalPort)
	_tunnel := tunnel.NewTunnelClient(
		tunnel.Forward,
		cfg.Protocol, fmt.Sprintf("%s:%d", cfg.Tunnel.Host, cfg.Tunnel.Port),
		cfg.Tunnel.CaCertPath, cfg.Tunnel.ServerName,
	)

	interceptor.OnTCPConnected = _tunnel.TunnelTCP
	interceptor.OnHTTPConnected = _tunnel.TunnelHTTP

	return &VPNClient{
		Protocol:    cfg.Protocol,
		Interceptor: interceptor,
		Tunnel:      _tunnel,
	}
}

func (c *VPNClient) Start() {
	go c.Interceptor.Start()
	go c.Tunnel.Start()

	select {}
}
