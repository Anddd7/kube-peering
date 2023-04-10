package tunnel

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/felixge/tcpkeepalive"
	"github.com/kube-peering/internal/pkg"
	"golang.org/x/net/http2"
)

func (t *tunnelServer) startHTTP() {
	if t.mode == pkg.Forward {
		server := http.Server{
			Addr:      t.localAddr(),
			TLSConfig: t.tlsConfig,
		}

		http.HandleFunc("/", t.onHTTPTunnel)

		t.logger.Infof("tunnel is listening on %s", t.localAddr())

		if err := server.ListenAndServeTLS("", ""); err != nil {
			log.Fatal(err)
		}
	}

	if t.mode == pkg.Reverse {
		tcpAddr, err := net.ResolveTCPAddr("tcp", t.localAddr())
		if err != nil {
			t.logger.Panicln(err)
		}

		listener, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			t.logger.Panicln(err)
		}
		defer listener.Close()

		t.logger.Infof("tunnel is listening on %s", t.localAddr())

		conn, err := listener.AcceptTCP()
		if err != nil {
			t.logger.Panicln(err)
		}

		tr := &http2.Transport{}
		tlsServerConn := tls.Server(conn, t.tlsConfig)
		h2Conn, err := tr.NewClientConn(tlsServerConn)
		if err != nil {
			// tls verification failed
			log.Fatal(err)
		}

		t.logger.Infof("tunnel is connected to %v", tlsServerConn.RemoteAddr().String())

		t.tlsConn = tlsServerConn
		t.clientConn = h2Conn

		tcpkeepalive.SetKeepAlive(t.tlsConn, 15*time.Minute, 3, 5*time.Second)
	}
}

func (t *tunnelServer) TunnelHTTP(w http.ResponseWriter, r *http.Request) {
	req := t.tunnelRequest(r)
	t.logger.Infof("tunnel request from [%s]%s to [%s]%s",
		r.RemoteAddr, r.URL.Path,
		t.tlsConn.RemoteAddr().String(), req.URL.Path,
	)

	resp, err := t.clientConn.RoundTrip(req)
	if err != nil {
		t.logger.Error(err)
		return
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		t.logger.Error(err)
	}
}

func (t *tunnelServer) tunnelRequest(r *http.Request) *http.Request {
	req := r.Clone(t.ctx)

	// clientHost := "localhost" //conn.RemoteAddr().(*net.TCPAddr).IP.String()
	// clientPort := conn.RemoteAddr().(*net.TCPAddr).Port

	// redirect to tunnel server
	req.URL = &url.URL{
		Scheme: "https",
		Path:   r.URL.String(),
	}

	// push the tunnel specific headers
	pushTunnelHeaders(req, r.Host)

	return req
}

func (t *tunnelServer) SetOnHTTPTunnel(fn http.HandlerFunc) {
	t.onHTTPTunnel = func(w http.ResponseWriter, r *http.Request) {
		t.logger.Infof("on http tunnel %s => %s\n", r.RemoteAddr, r.URL.String())

		req := r.Clone(t.ctx)

		// pop the tunnel specific headers
		_ = popTunnelHeaders(req)

		fn(w, req)
	}
}
