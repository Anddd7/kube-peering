package transit

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
	protocol        string
	port            int
	OnTCPConnected  func(conn *net.TCPConn)
	OnHTTPConnected http.HandlerFunc
}

func NewInterceptor(protocol string, port int) *Interceptor {
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
	if t.protocol == "tcp" {
		if t.OnTCPConnected == nil {
			t.logger.Panicln("OnTCPConnected is nil")
			panic("OnTCPConnected is needed for tcp protocol")
		}
		t.startTCP(t.OnTCPConnected)
	}

	if t.protocol == "http" {
		if t.OnHTTPConnected == nil {
			t.logger.Panicln("OnHTTPConnected is nil")
			panic("OnHTTPConnected is needed for http protocol")
		}
		t.startHttp(t.OnHTTPConnected)
	}
}

func (t *Interceptor) startTCP(onConnected func(conn *net.TCPConn)) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", t.port))
	if err != nil {
		t.logger.Panicln(err)
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		t.logger.Panicln(err)
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			t.logger.Error("", err)
			continue
		}
		go onConnected(conn)
	}
}

func (t *Interceptor) startHttp(onConnected http.HandlerFunc) {
	http2.ConfigureServer(&http.Server{}, &http2.Server{})
	http.ListenAndServe(fmt.Sprintf(":%d", t.port), onConnected)
}
