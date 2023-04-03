package transit

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

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

func NewTunnelClient(protocol, remoteAddr, caCertPath, serverName string) *TunnelClient {
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
		conn, err := tls.Dial("tcp", t.remoteAddr, t.tlsConfig)
		if err != nil {
			logger.Z.Errorf("failed to connect to %s: %v", t.remoteAddr, err)
			return
		}

		t.tlsConn = conn
		if t.onTCPTunnelIn != nil {
			t.onTCPTunnelIn(t.tlsConn)
		}
	}

	if t.protocol == "http" {
		// TODO
	}
}

func (t *TunnelClient) TunnelTCPOut(from *net.TCPConn) {
	for i := 0; i < 3; i++ {
		if t.tlsConn != nil {
			break
		}
		t.logger.Warnf("tunnel connection is nil, try to reconnect")
		time.Sleep(5 * time.Second)
	}

	if t.tlsConn == nil {
		t.logger.Panicln("tunnel connection is not ready")
		return
	}

	Pipe(t.logger, from, t.tlsConn)
}

func (t *TunnelClient) TunnelHttpOut(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (t *TunnelClient) SetOnTCPTunnelIn(fn func(conn *tls.Conn)) {
	t.onTCPTunnelIn = fn
}
