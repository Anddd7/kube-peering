package transit

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

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

func NewTunnelServer(protocol string, port int, serverCertPath, serverKeyPath, serverName string) *TunnelServer {
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
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", t.port))
	if err != nil {
		t.logger.Panicln(err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		t.logger.Panicln(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			t.logger.Error("", err)
			continue
		}
		go t.newConnection(conn)
	}
}

func (t *TunnelServer) newConnection(conn *net.TCPConn) {
	t.logger.Infof("new connection from %s", conn.RemoteAddr().String())

	tlsConn := tls.Server(conn, t.tlsConfig)

	t.tlsConn = tlsConn
	if t.protocol == "tcp" {
		if t.onTCPTunnelIn != nil {
			t.logger.Infof("running tcp tunnel with %s, wait for data ...", conn.RemoteAddr().String())

			t.onTCPTunnelIn(t.tlsConn)
		}
	}

	if t.protocol == "http" {
		// TODO
	}
}

func (t *TunnelServer) TunnelTCPOut(from *net.TCPConn) {
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

func (t *TunnelServer) SetOnTCPTunnelIn(fn func(conn *tls.Conn)) {
	t.onTCPTunnelIn = fn
}
func (t *TunnelServer) SetOnHttpTunnelIn(fn http.HandlerFunc) {
	t.onHttpTunnelIn = fn
}
