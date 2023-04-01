package connectors

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/felixge/tcpkeepalive"
	"github.com/kube-peering/internal/pkg/io"
	"github.com/kube-peering/internal/pkg/logger"
	"github.com/kube-peering/internal/pkg/model"
	"golang.org/x/net/http2"
)

/*
client (conn) <- TunnelServer <- requestChan
client (conn) -> TunnelServer -> responseChan
*/
func NewTunnelServer(
	ctx context.Context,
	cfg model.Tunnel,
	requestChan chan []byte,
	responseChan chan []byte,
) *TCPInterceptor {
	return &TCPInterceptor{
		ctx:       context.WithValue(ctx, keyComponentID, cfg.Name),
		mutex:     sync.Mutex{},
		wg:        sync.WaitGroup{},
		address:   cfg.Address(),
		readInto:  responseChan,
		writeFrom: requestChan,
	}
}

/*
responseChan -> TunnelClient -> server (conn)
requestChan <- TunnelClient <- server (conn)
*/
func NewTunnelClient(
	ctx context.Context,
	cfg model.Tunnel,
	requestChan chan []byte,
	responseChan chan []byte,
) *TCPForwarder {
	return &TCPForwarder{
		ctx:          context.WithValue(ctx, keyComponentID, cfg.Name),
		address:      cfg.Address(),
		forwardChan:  responseChan,
		backwordChan: requestChan,
	}
}

type TunnelServer struct {
	ctx        context.Context
	mutex      sync.Mutex
	ClientConn *http2.ClientConn
	address    string
	tlsConfig  *tls.Config
}

func NewHttp2TunnelServer(ctx context.Context, cfg model.Tunnel) *TunnelServer {
	crt, err := ioutil.ReadFile(cfg.ServerCertPath)
	if err != nil {
		log.Fatal(err)
	}

	key, err := ioutil.ReadFile(cfg.ServerKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	cert, err := tls.X509KeyPair(crt, key)
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   cfg.ServerName,
	}

	return &TunnelServer{
		ctx:       context.WithValue(ctx, keyComponentID, cfg.Name),
		address:   cfg.Address(),
		tlsConfig: tlsConfig,
	}
}

func (t *TunnelServer) Run() {
	io.StartTCPServer(
		t.address,
		func(s string) {
			logger.Z.Infof("[%s] server is started on %s", t.name(), s)
		},
		t.newConnection,
	)
}

func (t *TunnelServer) newConnection(conn net.Conn) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	logger.Z.Infof("[%s] new connection from %s", t.name(), conn.RemoteAddr().String())

	tlsConn := tls.Server(conn, t.tlsConfig)
	tr := &http2.Transport{}
	clientConn, err := tr.NewClientConn(tlsConn)

	if err != nil {
		logger.Z.Errorf("[%s] failed to create http2 connection: %v", t.name(), err)
		return
	}

	port := t.address[len(t.address)-4:]
	url := "https://localhost:" + fmt.Sprint(port)
	req, err := http.NewRequest(http.MethodConnect, url, nil)
	if err != nil {
		logger.Z.Errorf("[%s] failed to create initial request: %v", t.name(), err)
	}
	resp, err := clientConn.RoundTrip(req)
	if err != nil {
		logger.Z.Errorf("[%s] failed to send initial request: %v", t.name(), err)
	}

	if resp.StatusCode != http.StatusOK {
		clientConn.Close()
		return
	}

	t.ClientConn = clientConn

	logger.Z.Infof("[%s] http2 connection is established", t.name())
}

func (t *TunnelServer) name() string {
	return t.ctx.Value(keyComponentID).(string)
}

type TunnelClient struct {
	ctx       context.Context
	address   string
	tlsConfig *tls.Config
}

func NewHttp2TunnelClient(ctx context.Context, cfg model.Tunnel) *TunnelClient {
	crt, err := ioutil.ReadFile(cfg.CACertPath)
	if err != nil {
		log.Fatal(err)
	}

	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(crt)

	tlsConfig := &tls.Config{
		RootCAs:            rootCAs,
		InsecureSkipVerify: false,
		ServerName:         cfg.ServerName,
	}

	return &TunnelClient{
		ctx:       context.WithValue(ctx, keyComponentID, cfg.Name),
		address:   cfg.Address(),
		tlsConfig: tlsConfig,
	}
}

func (t *TunnelClient) Run() {
	conn, err := tls.Dial("tcp", t.address, t.tlsConfig)
	if err != nil {
		logger.Z.Errorf("[%s] failed to connect to %s: %v", t.name(), t.address, err)
		return
	}

	tcpkeepalive.SetKeepAlive(conn, 15*time.Minute, 3, 5*time.Second)

	h2s := &http2.Server{}
	h2s.ServeConn(conn, &http2.ServeConnOpts{
		Handler: http.HandlerFunc(t.proxyHttp),
	})
}

func (t *TunnelClient) proxyHttp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		t.handleHandshake(w, r)
		return
	}

	body, _ := ioutil.ReadAll(r.Body)
	msg := string(body)

	logger.Z.Infof("proto: %v", r.Proto)
	logger.Z.Infof("url: %v", r.URL)
	logger.Z.Infof("method: %v", r.Method)
	logger.Z.Infof("headers: %v", r.Header)
	logger.Z.Infof("body: %v", msg)

	if msg == "PING 5" {
		w.Write([]byte("You are luck boy!"))
		return
	}

	w.Write([]byte(
		fmt.Sprintf(
			"Reverse connection from client %s => %s, Protocol: %s",
			r.RemoteAddr, t.address, r.Proto,
		),
	))
}

func (t *TunnelClient) handleHandshake(w http.ResponseWriter, r *http.Request) {
	logger.Z.Infof("[%s] handshake with %s", t.name(), r.RemoteAddr)

	w.WriteHeader(http.StatusOK)
}

func (t *TunnelClient) name() string {
	return t.ctx.Value(keyComponentID).(string)
}
