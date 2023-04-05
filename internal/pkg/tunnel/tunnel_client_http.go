package tunnel

import (
	"net/http"
)

func (t *TunnelClient) startHttp() {
	// TODO
	// create http2 connection
	// wrap the income data into the customized http request
	// send the request to the server
	// unwrap the request and forward to the target application
	// send back the response similar like the request
}

func (t *TunnelClient) TunnelHttpOut(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (t *TunnelClient) SetOnHttpTunnelIn(fn http.HandlerFunc) {
	t.onHttpTunnelIn = fn
}
