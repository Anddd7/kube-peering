package tunnel

import (
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/http2"
)

func (t *TunnelClient) startHTTP() {
	// http2 to multiplex multiple requests over a single connection
	tr := &http2.Transport{
		TLSClientConfig:    t.tlsConfig,
		DisableCompression: true,
		AllowHTTP:          false,
	}
	client := &http.Client{
		Transport: tr,
	}

	t.httpClient = client
}

func (t *TunnelClient) TunnelHTTP(w http.ResponseWriter, r *http.Request) {
	req := t.tunnelRequest(r)
	t.logger.Infof("tunnel request from [%s]%s to [%s]%s",
		r.RemoteAddr, r.URL.Path,
		t.remoteAddr, req.URL.Path,
	)

	resp, err := t.httpClient.Do(req)
	if err != nil {
		t.logger.Error(err)
		return
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		t.logger.Error(err)
	}
}

func (t *TunnelClient) tunnelRequest(r *http.Request) *http.Request {
	req := r.Clone(t.ctx)

	// redirect to tunnel server
	req.RequestURI = ""
	req.URL = &url.URL{
		Scheme: "https",
		Host:   t.remoteAddr,
		Path:   r.URL.String(),
	}

	// push the tunnel specific headers
	pushTunnelHeaders(req, r.Host)

	return req
}

func (t *TunnelClient) SetOnHTTPTunnel(fn http.HandlerFunc) {
	// TODO
}
