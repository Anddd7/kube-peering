package connectors

import (
	"github.com/kube-peering/internal/pkg"
	"github.com/kube-peering/internal/pkg/tunnel"
)

type VPNConfig struct {
	Tunnel     pkg.TunnelConfig
	Protocol   pkg.Protocol
	LocalPort  int
	RemoteHost string
	RemotePort int
}

type VPNServer struct {
	Protocol  pkg.Protocol
	Forwarder *pkg.Forwarder
	Tunnel    pkg.Tunnel
}

func NewVPNServer(cfg VPNConfig) *VPNServer {
	_tunnel := tunnel.NewTunnelServer(
		pkg.Forward,
		cfg.Protocol, cfg.Tunnel.Port,
		cfg.Tunnel.ServerCertPath, cfg.Tunnel.ServerKeyPath, cfg.Tunnel.ServerName,
	)
	forwarder := pkg.NewForwarder(cfg.Protocol, cfg.RemoteHost, cfg.RemotePort)

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
	Tunnel      pkg.Tunnel
}

func NewVPNClient(cfg VPNConfig) *VPNClient {
	interceptor := pkg.NewInterceptor(cfg.Protocol, cfg.LocalPort)
	_tunnel := tunnel.NewTunnelClient(
		pkg.Forward,
		cfg.Protocol, cfg.Tunnel.Host, cfg.Tunnel.Port,
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
