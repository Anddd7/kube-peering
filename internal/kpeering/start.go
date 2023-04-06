package kpeering

import (
	"github.com/kube-peering/internal/pkg/connectors"
	"go.uber.org/zap"
)

type Kpeering struct {
	Logger    *zap.SugaredLogger
	VPNConfig connectors.VPNConfig
}

func (svr *Kpeering) Start() {
	// TODO pass the ctx and logger, reuse the base logger in sub components
	reverseVPNClient := connectors.NewReverseVPNServer(svr.VPNConfig)

	reverseVPNClient.Start()
}
