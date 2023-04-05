package tunnel

import (
	"net/http"
)

func (t *TunnelServer) startHTTP() {
	// TODO
}

func (t *TunnelServer) TunnelHTTPOut(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (t *TunnelServer) SetOnHTTPTunnelIn(fn http.HandlerFunc) {
	t.onHTTPTunnelIn = fn
}
