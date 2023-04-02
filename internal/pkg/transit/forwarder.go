package transit

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/kube-peering/internal/pkg/logger"
	"go.uber.org/zap"
)

type Forwarder struct {
	ctx          context.Context
	logger       *zap.SugaredLogger
	protocol     string
	remoteAddr   string
	reverseProxy *httputil.ReverseProxy
}

func NewFowarder(protocol string, remoteAddr string) *Forwarder {
	_logger := logger.CreateLocalLogger().With(
		"component", "proxy",
		"protocol", protocol,
	)

	var reverseProxy *httputil.ReverseProxy
	if protocol == "http" {
		remoteUrl, err := url.Parse(remoteAddr)
		if err != nil {
			_logger.Panicln(err)
		}
		reverseProxy = httputil.NewSingleHostReverseProxy(remoteUrl)
		reverseProxy.ModifyResponse = func(resp *http.Response) error {
			resp.Header.Set("X-Proxy-Server", "kube-peering")
			return nil
		}
	}

	return &Forwarder{
		ctx:          context.TODO(),
		logger:       _logger,
		protocol:     protocol,
		remoteAddr:   remoteAddr,
		reverseProxy: reverseProxy,
	}
}

func (t *Forwarder) ForwardTCP(from *net.TCPConn) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", t.remoteAddr)
	if err != nil {
		t.logger.Panicln(err)
	}

	to, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		t.logger.Panicln(err)
	}
	defer from.Close()
	defer to.Close()

	t.pipe(from, to)
}

func (t *Forwarder) pipe(from *net.TCPConn, to *net.TCPConn) {
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

func (t *Forwarder) ForwardHttp(w http.ResponseWriter, r *http.Request) {
	t.reverseProxy.ServeHTTP(w, r)
}
