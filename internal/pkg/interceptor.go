package pkg

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/kube-peering/internal/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
)

type Interceptor struct {
	ctx             context.Context
	logger          *zap.SugaredLogger
	protocol        Protocol
	port            int
	OnTCPConnected  func(conn PipeConn)
	OnHTTPConnected http.HandlerFunc
}

func NewInterceptor(protocol Protocol, port int) *Interceptor {
	return &Interceptor{
		ctx: context.TODO(),
		logger: logger.CreateLocalLogger().With(
			"component", "interceptor",
			"protocol", protocol,
		),
		protocol: protocol,
		port:     port,
	}
}

func (t *Interceptor) Start() {
	if t.protocol == TCP {
		if t.OnTCPConnected == nil {
			t.logger.Panicln("OnTCPConnected is nil")
		}
		t.startTCP(t.OnTCPConnected)
	}

	if t.protocol == HTTP {
		if t.OnHTTPConnected == nil {
			t.logger.Panicln("OnHTTPConnected is nil")
		}
		t.startHTTP(t.OnHTTPConnected)
	}
}

func (t *Interceptor) localAddr() string {
	return fmt.Sprintf(":%d", t.port)
}

func (t *Interceptor) startTCP(onConnected func(conn PipeConn)) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", t.localAddr())
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
		go onConnected(conn)
	}
}

func (t *Interceptor) startHTTP(onConnected http.HandlerFunc) {
	http2.ConfigureServer(&http.Server{}, &http2.Server{})
	http.ListenAndServe(t.localAddr(), onConnected)
}
