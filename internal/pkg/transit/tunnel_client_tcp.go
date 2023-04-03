package transit

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/kube-peering/internal/pkg/logger"
)

func (t *TunnelClient) startTCP() {
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

func (t *TunnelClient) SetOnTCPTunnelIn(fn func(conn *tls.Conn)) {
	t.onTCPTunnelIn = fn
}
