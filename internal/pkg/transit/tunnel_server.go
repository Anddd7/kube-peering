package transit

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/kube-peering/internal/pkg/config"
	"github.com/kube-peering/internal/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
)

// TunnelServer is a server that listens for incoming tunnel connections
type TunnelServer struct {
	ctx            context.Context
	logger         *zap.SugaredLogger
	protocol       string
	port           int
	tlsConfig      *tls.Config
	tlsConn        *tls.Conn
	onTCPTunnelIn  func(conn *tls.Conn)
	clientConn     *http2.ClientConn
	onHttpTunnelIn http.HandlerFunc
}

func NewTunnelServer(protocol string, port int, serverCertPath, serverKeyPath, serverName string) Tunnel {
	_logger := logger.CreateLocalLogger().With(
		"component", "tunnel",
		"mode", "server",
	)
	tlsConfig, err := config.LoadServerTlsConfig(serverCertPath, serverKeyPath, serverName)
	if err != nil {
		_logger.Panicln(err)
	}

	return &TunnelServer{
		ctx:       context.TODO(),
		logger:    _logger,
		protocol:  protocol,
		port:      port,
		tlsConfig: tlsConfig,
	}
}

func (t *TunnelServer) Start() {
	if t.protocol == "tcp" {
		t.startTCP()
	}
	if t.protocol == "http" {
		t.startHttp()
	}
}
