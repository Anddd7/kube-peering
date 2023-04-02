package transit

import (
	"context"

	"github.com/kube-peering/internal/pkg/logger"
	"go.uber.org/zap"
)

// TunnelServer is a server that listens for incoming tunnel connections
type TunnelServer struct {
	ctx    context.Context
	logger *zap.SugaredLogger
	port   int
}

func NewTunnelServer(port int) *TunnelServer {
	return &TunnelServer{
		ctx: context.TODO(),
		logger: logger.CreateLocalLogger().With(
			"component", "tunnel",
			"mode", "server",
		),
		port: port,
	}
}
