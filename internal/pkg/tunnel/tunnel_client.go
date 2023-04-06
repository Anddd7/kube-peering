package tunnel

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/kube-peering/internal/pkg/config"
	"github.com/kube-peering/internal/pkg/logger"
	"github.com/kube-peering/internal/pkg/util"
	"go.uber.org/zap"
)

// TunnelClient is a client that connects to a tunnel server
type tunnelClient struct {
	ctx          context.Context
	logger       *zap.SugaredLogger
	protocol     string
	remoteAddr   string
	tlsConfig    *tls.Config
	mode         TunnelMode
	tlsConn      *tls.Conn
	onTCPTunnel  func(conn util.PipeConn)
	httpClient   *http.Client
	onHTTPTunnel http.HandlerFunc
}

func NewTunnelClient(mode TunnelMode, protocol, remoteAddr, caCertPath, serverName string) Tunnel {
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
	if t.protocol == "tcp" {
		t.startTCP()
	}

	if t.protocol == "http" {
		t.startHTTP()
	}
}
