package tunnel

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/kube-peering/internal/pkg"
)

func (t *tunnelServer) startTCP() {
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
			t.logger.Error(err)
			continue
		}
		go t.newConnection(conn)
	}
}

func (t *tunnelServer) newConnection(conn *net.TCPConn) {
	t.logger.Infof("new connection from %s", conn.RemoteAddr().String())

	t.tlsConn = tls.Server(conn, t.tlsConfig)
	if t.onTCPTunnel != nil {
		t.logger.Infof("running tcp tunnel with %s, wait for data ...", conn.RemoteAddr().String())

		t.onTCPTunnel(t.tlsConn)
	}
}

func (t *tunnelServer) TunnelTCP(from pkg.PipeConn) {
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

func (t *tunnelServer) SetOnTCPTunnel(fn func(conn pkg.PipeConn)) {
	t.onTCPTunnel = fn
}
