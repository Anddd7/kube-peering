package tunnel

import (
	"fmt"
	"log"
	"net/http"
)

func (t *TunnelServer) startHTTP() {
	server := http.Server{
		Addr:      fmt.Sprintf(":%d", t.port),
		TLSConfig: t.tlsConfig,
	}

	http.HandleFunc("/", t.onHTTPTunnel)

	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatal(err)
	}
}

func (t *TunnelServer) TunnelHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (t *TunnelServer) SetOnHTTPTunnel(fn http.HandlerFunc) {
	t.onHTTPTunnel = func(w http.ResponseWriter, r *http.Request) {
		t.logger.Infof("on http tunnel %s => %s\n", r.RemoteAddr, r.URL.String())

		req := r.Clone(t.ctx)

		// pop the tunnel specific headers
		_ = popTunnelHeaders(req)

		fn(w, req)
	}
}

func pushTunnelHeaders(req *http.Request, host string) {
	req.Header.Set("X-Forwarded-Host", host)
}

func popTunnelHeaders(req *http.Request) string {
	host := req.Header.Get("X-Forwarded-Host")

	defer func(r *http.Request) {
		r.Header.Del("X-Forwarded-Host")
	}(req)

	return host
}
