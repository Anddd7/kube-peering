package pkg

import (
	"context"
	"fmt"
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
	mutex        sync.Mutex
	protocol     Protocol
	targetHost   string
	targetPort   int
	reverseProxy *httputil.ReverseProxy
}

func NewForwarder(protocol Protocol, targetHost string, targetPort int) *Forwarder {
	_logger := logger.CreateLocalLogger().With(
		"component", "forwarder",
		"protocol", protocol,
	)

	return &Forwarder{
		ctx:        context.TODO(),
		logger:     _logger,
		protocol:   protocol,
		targetHost: targetHost,
		targetPort: targetPort,
	}
}

func (t *Forwarder) remoteAddr() string {
	return fmt.Sprintf("%s:%d", t.targetHost, t.targetPort)
}

func (t *Forwarder) ForwardTCP(from PipeConn) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", t.remoteAddr())
	if err != nil {
		t.logger.Panicln(err)
	}

	to, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		t.logger.Panicln(err)
	}

	Pipe(t.logger, from, to)
}

func (t *Forwarder) initReverseProxy() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	targetUrl, err := url.Parse(fmt.Sprintf("%s://%s", t.protocol, t.remoteAddr()))
	if err != nil {
		t.logger.Panicln(err)
	}
	t.reverseProxy = httputil.NewSingleHostReverseProxy(targetUrl)
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
