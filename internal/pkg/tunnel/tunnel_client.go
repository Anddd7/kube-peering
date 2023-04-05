package tunnel

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/kube-peering/internal/pkg"
	"github.com/kube-peering/internal/pkg/config"
	"github.com/kube-peering/internal/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
)

// TunnelClient is a client that connects to a tunnel server
type TunnelClient struct {
	ctx            context.Context
	logger         *zap.SugaredLogger
	protocol       string
	remoteAddr     string
	tlsConfig      *tls.Config
	tlsConn        *tls.Conn
	onTCPTunnelIn  func(conn *tls.Conn)
	clientConn     *http2.ClientConn
	onHttpTunnelIn http.HandlerFunc
}

func NewTunnelClient(protocol, remoteAddr, caCertPath, serverName string) pkg.Tunnel {
	_logger := logger.CreateLocalLogger().With(
		"component", "tunnel",
		"mode", "client",
	)
	tlsConfig, err := config.LoadClientTlsConfig(caCertPath, serverName)
	if err != nil {
		_logger.Panicln(err)
	}

	return &TunnelClient{
		ctx:        context.TODO(),
		logger:     _logger,
		protocol:   protocol,
		remoteAddr: remoteAddr,
		tlsConfig:  tlsConfig,
	}
}

func (t *TunnelClient) Start() {
	if t.protocol == "tcp" {
		t.startTCP()
	}

	if t.protocol == "http" {
		t.startHttp()
	}
}
