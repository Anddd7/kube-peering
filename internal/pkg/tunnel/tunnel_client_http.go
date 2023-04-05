package tunnel

import (
	"net/http"
)

func (t *TunnelClient) startHTTP() {
	// TODO
	// create http2 connection
	// wrap the income data into the customized http request
	// send the request to the server
	// unwrap the request and forward to the target application
	// send back the response similar like the request
}

func (t *TunnelClient) TunnelHTTPOut(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (t *TunnelClient) SetOnHTTPTunnelIn(fn http.HandlerFunc) {
	t.onHTTPTunnelIn = fn
}
