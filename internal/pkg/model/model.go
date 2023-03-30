package model

var DefaultTunnel = CreateTunnel("localhost", 10022)

type Interceptor struct {
	Endpoint
}

type Tunnel struct {
	Endpoint
}

type Forwarder struct {
	Endpoint
}

func CreateInterceptor(protocol string, port int) Interceptor {
	return Interceptor{
		Endpoint: Endpoint{
			Name:       "interceptor",
			Protocol:   protocol,
			Host:       "localhost",
			ListenPort: port,
		},
	}
}

func CreateTunnel(host string, port int) Tunnel {
	return Tunnel{
		Endpoint: Endpoint{
			Name:       "tunnel",
			Protocol:   "tcp",
			Host:       host,
			ListenPort: port,
		},
	}
}

func CreateForwarder(host string, port int) Forwarder {
	return Forwarder{
		Endpoint: Endpoint{
			Host:       host,
			ListenPort: port,
		},
	}
}
