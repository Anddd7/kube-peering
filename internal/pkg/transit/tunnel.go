package transit

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"sync"

	"go.uber.org/zap"
)

type Tunnel interface {
	Start()
	SetOnTCPTunnelIn(func(conn *tls.Conn))
	SetOnHttpTunnelIn(http.HandlerFunc)
	TunnelTCPOut(from *net.TCPConn)
	TunnelHttpOut(w http.ResponseWriter, r *http.Request)
}

func Pipe(_logger *zap.SugaredLogger, from, to PipeConn) {
	defer from.Close()
	defer to.Close()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		_logger.Infof("transfer data from %s to %s", from.RemoteAddr().String(), to.RemoteAddr().String())
		io.Copy(from, to)
		wg.Done()
	}()

	go func() {
		_logger.Infof("transfer data back from %s to %s", to.RemoteAddr().String(), from.RemoteAddr().String())
		io.Copy(to, from)
		wg.Done()
	}()

	wg.Wait()
}

type PipeConn interface {
	io.Reader
	io.Writer
	io.Closer
	RemoteAddr() net.Addr
}
