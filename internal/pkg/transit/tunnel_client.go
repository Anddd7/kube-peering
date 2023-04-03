package transit

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"sync"
	"time"

	"github.com/kube-peering/internal/pkg/config"
	"github.com/kube-peering/internal/pkg/logger"
	"go.uber.org/zap"
)

// TunnelClient is a client that connects to a tunnel server
type TunnelClient struct {
	ctx        context.Context
	logger     *zap.SugaredLogger
	remoteAddr string
	tlsConfig  *tls.Config
	mutex      sync.Mutex
	tlsConn    *tls.Conn
}

func NewTunnelClient(remoteAddr, caCertPath, serverName string) *TunnelClient {
	_logger := logger.CreateLocalLogger().With(
		"component", "tunnel",
		"mode", "server",
	)
	tlsConfig, err := config.LoadClientTlsConfig(caCertPath, serverName)
	if err != nil {
		_logger.Panicln(err)
	}

	return &TunnelClient{
		ctx:        context.TODO(),
		logger:     _logger,
		remoteAddr: remoteAddr,
		tlsConfig:  tlsConfig,
	}
}

func (t *TunnelClient) Start() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	conn, err := tls.Dial("tcp", t.remoteAddr, t.tlsConfig)
	if err != nil {
		logger.Z.Errorf("failed to connect to %s: %v", t.remoteAddr, err)
		return
	}

	t.tlsConn = conn

	// tcpkeepalive.SetKeepAlive(conn, 15*time.Minute, 3, 5*time.Second)

	// h2s := &http2.Server{}
	// h2s.ServeConn(conn, &http2.ServeConnOpts{
	// 	Handler: http.HandlerFunc(t.proxyHttp),
	// })
}

func (t *TunnelClient) ForwardTls(conn *net.TCPConn) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for i := 0; i < 3; i++ {
		if t.tlsConn != nil {
			break
		}
		t.logger.Warnf("tls connection is nil, try to reconnect")
		time.Sleep(5 * time.Second)
	}

	t.pipe(conn, t.tlsConn)
}

func (t *TunnelClient) pipe(from *net.TCPConn, to *tls.Conn) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		t.logger.Infof("transfer data from %s to %s", from.RemoteAddr().String(), to.RemoteAddr().String())
		io.Copy(from, to)
		wg.Done()
	}()

	go func() {
		t.logger.Infof("transfer data back from %s to %s", to.RemoteAddr().String(), from.RemoteAddr().String())
		io.Copy(to, from)
		wg.Done()
	}()

	wg.Wait()
}
