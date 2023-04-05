package tunnel

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/http2"
)

func (t *TunnelClient) startHTTP() {
	client := &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig:    t.tlsConfig,
			DisableCompression: true,
			AllowHTTP:          false,
		},
	}

	t.httpClient = client

	// TODO
	// create http2 connection
	// wrap the income data into the customized http request
	// send the request to the server
	// unwrap the request and forward to the target application
	// send back the response similar like the request
}

func (t *TunnelClient) TunnelHTTP(w http.ResponseWriter, r *http.Request) {
	pr, pw := io.Pipe()
	req := t.tunnelRequest(r)
	resp, err := t.httpClient.Do(req)
	if err != nil {
		t.logger.Errorln(err)
		return
	}

	t.logger.Infof("content length: %d", r.ContentLength)
	t.logger.Infof("content length: %d", req.ContentLength)

	go io.Copy(pw, resp.Body)
	go io.Copy(w, pr)
}

func (t *TunnelClient) tunnelRequest(r *http.Request) *http.Request {
	req, _ := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://%s%s", t.remoteAddr, r.URL.String()),
		r.Body,
	)
	req.Header = r.Header.Clone()
	req.Header.Set("X-Forwarded-For", r.RemoteAddr)
	req.Header.Set("X-Forwarded-Proto", r.Proto)
	req.Header.Set("X-Forwarded-Method", r.Method)
	req.Header.Set("X-Forwarded-Scheme", r.URL.Scheme)
	req.Header.Set("X-Forwarded-Host", r.Host)
	req.Header.Set("X-Forwarded-Path", r.URL.String())

	return req
}

func (t *TunnelClient) SetOnHTTPTunnel(fn http.HandlerFunc) {
	t.onHTTPTunnel = fn
}
