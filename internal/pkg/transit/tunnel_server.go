package transit

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/kube-peering/internal/pkg/config"
	"github.com/kube-peering/internal/pkg/logger"
	"go.uber.org/zap"
)

// TunnelServer is a server that listens for incoming tunnel connections
type TunnelServer struct {
	ctx           context.Context
	logger        *zap.SugaredLogger
	port          int
	tlsConfig     *tls.Config
	tlsConn       *tls.Conn
	onTCPTunnelIn func(conn *tls.Conn)
}

func NewTunnelServer(port int, serverCertPath, serverKeyPath, serverName string) *TunnelServer {
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
	if t.onTCPTunnelIn != nil {
		t.OnTCPTunnelIn(t.tlsConn)
	}

	// tlsConn, clientConn, err := t.initClientConn(conn)
	// if err != nil {
	// 	t.logger.Errorf("failed to create http2 connection: %v", err)
	// 	return
	// }

	// t.clientConn = clientConn
}

func (t *TunnelServer) TunnelOut(from *net.TCPConn) {
	for i := 0; i < 3; i++ {
		if t.tlsConn != nil {
			break
		}
		t.logger.Warnf("tls connection is nil, try to reconnect")
		time.Sleep(5 * time.Second)
	}

	Pipe(t.logger, from, t.tlsConn)
}

func (t *TunnelServer) SetOnTCPTunnelIn(fn func(conn *tls.Conn)) {
	t.onTCPTunnelIn = fn
}
func (t *TunnelServer) OnTCPTunnelIn(from *tls.Conn) {
	t.onTCPTunnelIn(from)
}

// func (t *TunnelServer) initClientConn(conn *net.TCPConn) (*tls.Conn, *http2.ClientConn, error) {
// 	tlsConn := tls.Server(conn, t.tlsConfig)
// 	tr := &http2.Transport{}
// 	clientConn, err := tr.NewClientConn(tlsConn)
// 	if err != nil {
// 		conn.Close()
// 		return nil, nil, err
// 	}

// 	return tlsConn, clientConn, nil
// }
