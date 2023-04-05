package tunnel

import (
	"net/http"
)

func (t *TunnelServer) startHttp() {
	// TODO
}

func (t *TunnelServer) TunnelHttpOut(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (t *TunnelServer) SetOnHttpTunnelIn(fn http.HandlerFunc) {
	t.onHttpTunnelIn = fn
}
