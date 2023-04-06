package tunnel

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/kube-peering/internal/pkg"
	"github.com/kube-peering/internal/pkg/config"
	"github.com/kube-peering/internal/pkg/logger"
	"go.uber.org/zap"
)

// TunnelClient is a client that connects to a tunnel server
type tunnelClient struct {
	ctx          context.Context
	logger       *zap.SugaredLogger
	protocol     pkg.Protocol
	remoteAddr   string
	tlsConfig    *tls.Config
	mode         pkg.TunnelMode
	tlsConn      *tls.Conn
	onTCPTunnel  func(conn pkg.PipeConn)
	httpClient   *http.Client
	onHTTPTunnel http.HandlerFunc
}

func NewTunnelClient(mode pkg.TunnelMode, protocol pkg.Protocol, remoteAddr, caCertPath, serverName string) pkg.Tunnel {
	_logger := logger.CreateLocalLogger().With(
		"component", "tunnel",
		"type", "client",
		"mode", mode.String(),
		"protocol", protocol,
	)
	tlsConfig, err := config.LoadClientTlsConfig(caCertPath, serverName)
	if err != nil {
		_logger.Panicln(err)
	}

	return &tunnelClient{
		ctx:        context.TODO(),
		logger:     _logger,
		protocol:   protocol,
		remoteAddr: remoteAddr,
		tlsConfig:  tlsConfig,
		mode:       mode,
	}
}

func (t *tunnelClient) Start() {
	switch t.protocol {
	case pkg.TCP:
		t.startTCP()
	case pkg.HTTP:
		t.startHTTP()
	}
}
