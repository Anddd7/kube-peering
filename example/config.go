package http2

import "github.com/kube-peering/internal/pkg/config"

var (
	AppTlsConfig, _    = config.LoadServerTlsConfig("../../certs_app/server.crt", "../../certs_app/server.key", "localhost")
	ClientTlsConfig, _ = config.LoadClientTlsConfig("../../certs_app/ca.crt", "localhost")
)

const (
	TunnelServerCert = "../../certs_tunnel/server.crt"
	TunnelServerKey  = "../../certs_tunnel/server.key"
	TunnelCaCert     = "../../certs_tunnel/ca.crt"
	TunnelServerName = "localhost"
)

var (
	TunnelServerTlsConfig, _ = config.LoadServerTlsConfig(TunnelServerCert, TunnelServerKey, TunnelServerName)
	TunnelClientTlsConfig, _ = config.LoadClientTlsConfig(TunnelCaCert, TunnelServerName)
)

const (
	AppPort = 8080
	AppAddr = "localhost:8080"

	AppHttpsPort = 8443
	AppHttpsAddr = "localhost:8443"

	TunnelPort = 10086
	TunnelHost = "localhost"
	TunnelAddr = "localhost:10086"

	ProxyPort = 10021
	ProxyAddr = "localhost:10021"

	VPNPort = 10022
	VPNAddr = "localhost:10022"
)
