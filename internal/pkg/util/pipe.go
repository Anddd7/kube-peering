package util

import (
	"io"
	"net"
	"sync"

	"go.uber.org/zap"
)

// interface to represent a net.Conn or net.TCPConn or tls.Conn
type PipeConn interface {
	io.Reader
	io.Writer
	io.Closer
	RemoteAddr() net.Addr
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
