package connectors

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kube-peering/internal/pkg"
	"github.com/kube-peering/internal/pkg/tunnel"
)

type PortForwardServer struct {
	Protocol    pkg.Protocol
	Interceptor *pkg.Interceptor
	Tunnel      tunnel.Tunnel
}

func NewPortFowardServer(cfg VPNConfig) *PortForwardServer {
	remotePort, _ := strconv.Atoi(strings.Split(cfg.RemoteAddr, ":")[1])

	interceptor := pkg.NewInterceptor(cfg.Protocol, remotePort)
	_tunnel := tunnel.NewTunnelServer(
		tunnel.Reverse,
		cfg.Protocol, cfg.Tunnel.Port,
		cfg.Tunnel.ServerCertPath, cfg.Tunnel.ServerKeyPath, cfg.Tunnel.ServerName,
	)
	interceptor.OnTCPConnected = _tunnel.TunnelTCP
	interceptor.OnHTTPConnected = _tunnel.TunnelHTTP

	return &PortForwardServer{
		Protocol:    cfg.Protocol,
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

func NewPortForwardClient(cfg VPNConfig) *PortForwardClient {
	localAddr := fmt.Sprintf("localhost:%d", cfg.LocalPort)

	_tunnel := tunnel.NewTunnelClient(
		tunnel.Forward,
		cfg.Protocol, fmt.Sprintf("%s:%d", cfg.Tunnel.Host, cfg.Tunnel.Port),
		cfg.Tunnel.CaCertPath, cfg.Tunnel.ServerName,
	)
	forwarder := pkg.NewForwarder(cfg.Protocol, localAddr)

	_tunnel.SetOnTCPTunnel(forwarder.ForwardTCP)
	_tunnel.SetOnHTTPTunnel(forwarder.ForwardHTTP)

	return &PortForwardClient{
		Protocol:  cfg.Protocol,
		Forwarder: forwarder,
		Tunnel:    _tunnel,
	}
}

func (c *PortForwardClient) Start() {
	c.Tunnel.Start()
}
