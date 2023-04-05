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
	t.onHTTPTunnel = fn
}
