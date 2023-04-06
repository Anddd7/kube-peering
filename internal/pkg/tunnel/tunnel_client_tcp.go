package tunnel

import (
	"crypto/tls"
	"time"

	"github.com/kube-peering/internal/pkg"
)

func (t *tunnelClient) startTCP() {
	conn, err := tls.Dial("tcp", t.remoteAddr, t.tlsConfig)
	if err != nil {
		t.logger.Errorf("failed to connect to %s: %v", t.remoteAddr, err)
		return
	}

	t.tlsConn = conn

	if t.onTCPTunnel != nil {
		t.onTCPTunnel(t.tlsConn)
	}
}

func (t *tunnelClient) TunnelTCP(from pkg.PipeConn) {
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

	pkg.Pipe(t.logger, from, t.tlsConn)
}

func (t *tunnelClient) SetOnTCPTunnel(fn func(conn pkg.PipeConn)) {
	t.onTCPTunnel = fn
}
