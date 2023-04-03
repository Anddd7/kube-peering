package http2

import "github.com/kube-peering/internal/pkg/config"

var AppTlsConfig, _ = config.LoadServerTlsConfig("../../certs_app/server.crt", "../../certs_app/server.key", "localhost")
var ClientTlsConfig, _ = config.LoadClientTlsConfig("../../certs_app/ca.crt", "localhost")

const (
	TunnelServerCert = "../../certs_tunnel/server.crt"
	TunnelServerKey  = "../../certs_tunnel/server.key"
	TunnelCaCert     = "../../certs_tunnel/ca.crt"
	TunnelServerName = "localhost"
)

var TunnelServerTlsConfig, _ = config.LoadServerTlsConfig(TunnelServerCert, TunnelServerKey, TunnelServerName)
var TunnelClientTlsConfig, _ = config.LoadClientTlsConfig(TunnelCaCert, TunnelServerName)
