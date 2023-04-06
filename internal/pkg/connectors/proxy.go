package connectors

import "github.com/kube-peering/internal/pkg"

type Proxy struct {
	Protocol    pkg.Protocol
	Interceptor *pkg.Interceptor
	Forwarder   *pkg.Forwarder
}

func NewProxy(protocol pkg.Protocol, localPort int, remoteAddr string) *Proxy {
	interceptor := pkg.NewInterceptor(protocol, localPort)
	forwarder := pkg.NewForwarder(protocol, remoteAddr)

	interceptor.OnTCPConnected = forwarder.ForwardTCP
	interceptor.OnHTTPConnected = forwarder.ForwardHTTP

	return &Proxy{
		Protocol:    protocol,
		Interceptor: interceptor,
		Forwarder:   forwarder,
	}
}

func (t *Proxy) Start() {
	t.Interceptor.Start()
}
