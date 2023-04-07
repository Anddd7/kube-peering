package connectors

import (
	"github.com/kube-peering/internal/pkg"
	"github.com/kube-peering/internal/pkg/tunnel"
)

type ReverseVPNServer struct {
	Protocol    pkg.Protocol
	Interceptor *pkg.Interceptor
	Tunnel      pkg.Tunnel
}

func NewReverseVPNServer(cfg VPNConfig) *ReverseVPNServer {
	interceptor := pkg.NewInterceptor(cfg.Protocol, cfg.RemotePort)
	_tunnel := tunnel.NewTunnelServer(
		pkg.Reverse,
		cfg.Protocol, cfg.Tunnel.Port,
		cfg.Tunnel.ServerCertPath, cfg.Tunnel.ServerKeyPath, cfg.Tunnel.ServerName,
	)
	interceptor.OnTCPConnected = _tunnel.TunnelTCP
	interceptor.OnHTTPConnected = _tunnel.TunnelHTTP

	return &ReverseVPNServer{
		Protocol:    cfg.Protocol,
		Interceptor: interceptor,
		Tunnel:      _tunnel,
	}
}

func (s *ReverseVPNServer) Start() {
	go s.Interceptor.Start()
	go s.Tunnel.Start()

	select {}
}

type ReverseVPNClient struct {
	Protocol  pkg.Protocol
	Forwarder *pkg.Forwarder
	Tunnel    pkg.Tunnel
}

func NewReverseVPNClient(cfg VPNConfig) *ReverseVPNClient {
	_tunnel := tunnel.NewTunnelClient(
		pkg.Forward,
		cfg.Protocol, cfg.Tunnel.Host, cfg.Tunnel.Port,
		cfg.Tunnel.CaCertPath, cfg.Tunnel.ServerName,
	)
	forwarder := pkg.NewForwarder(cfg.Protocol, "localhost", cfg.LocalPort)

	_tunnel.SetOnTCPTunnel(forwarder.ForwardTCP)
	_tunnel.SetOnHTTPTunnel(forwarder.ForwardHTTP)

	return &ReverseVPNClient{
		Protocol:  cfg.Protocol,
		Forwarder: forwarder,
		Tunnel:    _tunnel,
	}
}

func (c *ReverseVPNClient) Start() {
	c.Tunnel.Start()
}
