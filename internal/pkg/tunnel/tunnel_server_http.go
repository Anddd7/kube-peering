package tunnel

import (
	"net/http"
)

func (t *TunnelServer) startHTTP() {
	// TODO
}

func (t *TunnelServer) TunnelHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (t *TunnelServer) SetOnHTTPTunnel(fn http.HandlerFunc) {
	t.onHTTPTunnel = fn
}
