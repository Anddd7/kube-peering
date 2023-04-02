package http2

import "github.com/kube-peering/internal/pkg/config"

var appTlsConfig, _ = config.LoadServerTlsConfig("../../certs_app/server.crt", "../../certs_app/server.key", "localhost")
var clientTlsConfig, _ = config.LoadClientTlsConfig("../../certs_app/ca.crt", "localhost")

var tunnelServerTlsConfig, _ = config.LoadServerTlsConfig("../../certs_tunnel/server.crt", "../../certs_app/server.key", "localhost")
var tunnelClientTlsConfig, _ = config.LoadClientTlsConfig("../../certs_tunnel/ca.crt", "localhost")
