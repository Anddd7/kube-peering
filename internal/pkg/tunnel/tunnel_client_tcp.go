package tunnel

import (
	"crypto/tls"
	"time"

	"github.com/kube-peering/internal/pkg/logger"
	"github.com/kube-peering/internal/pkg/util"
)

func (t *TunnelClient) startTCP() {
	conn, err := tls.Dial("tcp", t.remoteAddr, t.tlsConfig)
	if err != nil {
		logger.Z.Errorf("failed to connect to %s: %v", t.remoteAddr, err)
		return
	}

	t.tlsConn = conn

	if t.onTCPTunnel != nil {
		t.onTCPTunnel(t.tlsConn)
	}
}

func (t *TunnelClient) TunnelTCP(from util.PipeConn) {
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

	util.Pipe(t.logger, from, t.tlsConn)
}

func (t *TunnelClient) SetOnTCPTunnel(fn func(conn util.PipeConn)) {
	t.onTCPTunnel = fn
}
