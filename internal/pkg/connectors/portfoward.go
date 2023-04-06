package connectors

import (
	"github.com/kube-peering/internal/pkg"
	"github.com/kube-peering/internal/pkg/tunnel"
)

type PortForwardServer struct {
	Protocol    pkg.Protocol
	Interceptor *pkg.Interceptor
	Tunnel      tunnel.Tunnel
}

func NewPortFowardServer(protocol pkg.Protocol, tunnelPort int, localPort int, serverCertPath, serverKeyPath, serverName string) *PortForwardServer {
	interceptor := pkg.NewInterceptor(protocol, localPort)
	_tunnel := tunnel.NewTunnelServer(tunnel.Reverse, protocol, tunnelPort, serverCertPath, serverKeyPath, serverName)

	interceptor.OnTCPConnected = _tunnel.TunnelTCP
	interceptor.OnHTTPConnected = _tunnel.TunnelHTTP

	return &PortForwardServer{
		Protocol:    protocol,
		Interceptor: interceptor,
		Tunnel:      _tunnel,
	}
}

func (s *PortForwardServer) Start() {
	go s.Interceptor.Start()
	go s.Tunnel.Start()

	select {}
}

type PortForwardClient struct {
	Protocol  pkg.Protocol
	Forwarder *pkg.Forwarder
	Tunnel    tunnel.Tunnel
}

func NewPortForwardClient(protocol pkg.Protocol, tunnelAddr string, localAddr string, caCertPath string, serverName string) *PortForwardClient {
	_tunnel := tunnel.NewTunnelClient(tunnel.Forward, protocol, tunnelAddr, caCertPath, serverName)
	forwarder := pkg.NewForwarder(protocol, localAddr)

	_tunnel.SetOnTCPTunnel(forwarder.ForwardTCP)
	_tunnel.SetOnHTTPTunnel(forwarder.ForwardHTTP)

	return &PortForwardClient{
		Protocol:  protocol,
		Forwarder: forwarder,
		Tunnel:    _tunnel,
	}
}

func (c *PortForwardClient) Start() {
	c.Tunnel.Start()
}
