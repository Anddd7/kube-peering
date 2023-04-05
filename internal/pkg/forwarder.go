package pkg

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/kube-peering/internal/pkg/logger"
	"github.com/kube-peering/internal/pkg/util"
	"go.uber.org/zap"
)

type Forwarder struct {
	ctx          context.Context
	logger       *zap.SugaredLogger
	mutex        sync.Mutex
	protocol     string
	remoteAddr   string
	reverseProxy *httputil.ReverseProxy
}

func NewFowarder(protocol string, remoteAddr string) *Forwarder {
	_logger := logger.CreateLocalLogger().With(
		"component", "forwarder",
		"protocol", protocol,
	)

	return &Forwarder{
		ctx:        context.TODO(),
		logger:     _logger,
		protocol:   protocol,
		remoteAddr: remoteAddr,
	}
}

func (t *Forwarder) ForwardTCP(from util.PipeConn) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", t.remoteAddr)
	if err != nil {
		t.logger.Panicln(err)
	}

	to, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		t.logger.Panicln(err)
	}

	util.Pipe(t.logger, from, to)
}

func (t *Forwarder) initReverseProxy() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	remoteUrl, err := url.Parse(t.remoteAddr)
	if err != nil {
		t.logger.Panicln(err)
	}
	t.reverseProxy = httputil.NewSingleHostReverseProxy(remoteUrl)
	t.reverseProxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("X-Proxy-Server", "kube-peering")
		return nil
	}
}

func (t *Forwarder) ForwardHTTP(w http.ResponseWriter, r *http.Request) {
	if t.reverseProxy == nil {
		t.initReverseProxy()
	}
	t.reverseProxy.ServeHTTP(w, r)
}
