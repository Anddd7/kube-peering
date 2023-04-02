package transit

import (
	"context"

	"github.com/kube-peering/internal/pkg/logger"
	"go.uber.org/zap"
)

// TunnelClient is a client that connects to a tunnel server
type TunnelClient struct {
	ctx        context.Context
	logger     *zap.SugaredLogger
	remoteAddr string
}

func NewTunnelClient(remoteAddr string) *TunnelClient {
	return &TunnelClient{
		ctx: context.TODO(),
		logger: logger.CreateLocalLogger().With(
			"component", "tunnel",
			"mode", "client",
		),
		remoteAddr: remoteAddr,
	}
}
