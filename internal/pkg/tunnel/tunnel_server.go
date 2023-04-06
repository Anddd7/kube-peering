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

// TunnelServer is a server that listens for incoming tunnel connections
type tunnelServer struct {
	ctx          context.Context
	logger       *zap.SugaredLogger
	protocol     pkg.Protocol
	port         int
	tlsConfig    *tls.Config
	mode         pkg.TunnelMode
	tlsConn      *tls.Conn
	onTCPTunnel  func(conn pkg.PipeConn)
	clientConn   *http2.ClientConn
	onHTTPTunnel http.HandlerFunc
}

func NewTunnelServer(mode pkg.TunnelMode, protocol pkg.Protocol, port int, serverCertPath, serverKeyPath, serverName string) pkg.Tunnel {
	_logger := logger.CreateLocalLogger().With(
		"component", "tunnel",
		"type", "server",
		"mode", mode.String(),
		"protocol", protocol,
	)
	tlsConfig, err := config.LoadServerTlsConfig(serverCertPath, serverKeyPath, serverName)
	if err != nil {
		_logger.Panicln(err)
	}

	return &tunnelServer{
		ctx:       context.TODO(),
		logger:    _logger,
		protocol:  protocol,
		port:      port,
		tlsConfig: tlsConfig,
		mode:      mode,
	}
}

func (t *tunnelServer) Start() {
	switch t.protocol {
	case pkg.TCP:
		t.startTCP()
	case pkg.HTTP:
		t.startHTTP()
	}
}

func pushTunnelHeaders(req *http.Request, host string) {
	req.Header.Set("X-Forwarded-Host", host)
}

func popTunnelHeaders(req *http.Request) string {
	host := req.Header.Get("X-Forwarded-Host")

	defer func(r *http.Request) {
		r.Header.Del("X-Forwarded-Host")
	}(req)

	return host
}
