package connectors

import (
	"github.com/kube-peering/internal/pkg"
	"github.com/kube-peering/internal/pkg/tunnel"
)

type VPNServer struct {
	Protocol  pkg.Protocol
	Forwarder *pkg.Forwarder
	Tunnel    tunnel.Tunnel
}

func NewVPNServer(protocol pkg.Protocol, tunnelPort int, remoteAddr string, serverCertPath, serverKeyPath, serverName string) *VPNServer {
	_tunnel := tunnel.NewTunnelServer(tunnel.Forward, protocol, tunnelPort, serverCertPath, serverKeyPath, serverName)
	forwarder := pkg.NewForwarder(protocol, remoteAddr)

	_tunnel.SetOnTCPTunnel(forwarder.ForwardTCP)
	_tunnel.SetOnHTTPTunnel(forwarder.ForwardHTTP)

	return &VPNServer{
		Protocol:  protocol,
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

func NewVPNClient(protocol pkg.Protocol, localPort int, tunnelAddr string, caCertPath string, serverName string) *VPNClient {
	interceptor := pkg.NewInterceptor(protocol, localPort)
	_tunnel := tunnel.NewTunnelClient(tunnel.Forward, protocol, tunnelAddr, caCertPath, serverName)

	interceptor.OnTCPConnected = _tunnel.TunnelTCP
	interceptor.OnHTTPConnected = _tunnel.TunnelHTTP

	return &VPNClient{
		Protocol:    protocol,
		Interceptor: interceptor,
		Tunnel:      _tunnel,
	}
}

func (c *VPNClient) Start() {
	go c.Interceptor.Start()
	go c.Tunnel.Start()

	select {}
}
