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

type Proxy struct {
	ctx          context.Context
	logger       *zap.SugaredLogger
	protocol     string
	remoteAddr   string
	reverseProxy *httputil.ReverseProxy
}

func NewProxy(protocol string, remoteAddr string) *Proxy {
	_logger := logger.CreateLocalLogger().With(
		"component", "proxy",
		"protocol", protocol,
	)

	var reverseProxy httputil.ReverseProxy
	if protocol == "http" {
		remoteUrl, err := url.Parse(remoteAddr)
		if err != nil {
			_logger.Panicln(err)
			panic(err)
		}
		reverseProxy := httputil.NewSingleHostReverseProxy(remoteUrl)
		reverseProxy.ModifyResponse = func(resp *http.Response) error {
			resp.Header.Set("X-Proxy-Server", "kube-peering")
			return nil
		}
	}

	return &Proxy{
		ctx:          context.TODO(),
		logger:       _logger,
		protocol:     protocol,
		remoteAddr:   remoteAddr,
		reverseProxy: &reverseProxy,
	}
}

func (t *Proxy) ProxyTCP(from *net.TCPConn) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", t.remoteAddr)
	if err != nil {
		panic(err)
	}

	to, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}
	defer from.Close()
	defer to.Close()

	t.pipe(from, to)
}

func (t *Proxy) pipe(from *net.TCPConn, to *net.TCPConn) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		t.logger.Infof("transfer data from %s to %s", from.RemoteAddr().String(), to.RemoteAddr().String())
		io.Copy(from, to)
		wg.Done()
	}()

	go func() {
		t.logger.Info("transfer data back from %s to %s", to.RemoteAddr().String(), from.RemoteAddr().String())
		io.Copy(to, from)
		wg.Done()
	}()

	wg.Wait()
}

func (t *Proxy) ProxyHttp(w http.ResponseWriter, r *http.Request) {
	t.reverseProxy.ServeHTTP(w, r)
}
