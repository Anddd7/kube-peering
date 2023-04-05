package tunnel

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/kube-peering/internal/pkg/util"
)

func (t *TunnelServer) startTCP() {
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

	t.tlsConn = tls.Server(conn, t.tlsConfig)
	if t.onTCPTunnelIn != nil {
		t.logger.Infof("running tcp tunnel with %s, wait for data ...", conn.RemoteAddr().String())

		t.onTCPTunnelIn(t.tlsConn)
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

	util.Pipe(t.logger, from, t.tlsConn)
}

func (t *TunnelServer) SetOnTCPTunnelIn(fn func(conn *tls.Conn)) {
	t.onTCPTunnelIn = fn
}
