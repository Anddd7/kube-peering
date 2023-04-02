package http2

import "github.com/kube-peering/internal/pkg/config"

var AppTlsConfig, _ = config.LoadServerTlsConfig("../../certs_app/server.crt", "../../certs_app/server.key", "localhost")
var ClientTlsConfig, _ = config.LoadClientTlsConfig("../../certs_app/ca.crt", "localhost")

var TunnelServerTlsConfig, _ = config.LoadServerTlsConfig("../../certs_tunnel/server.crt", "../../certs_app/server.key", "localhost")
var TunnelClientTlsConfig, _ = config.LoadClientTlsConfig("../../certs_tunnel/ca.crt", "localhost")
