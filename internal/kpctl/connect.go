package kpctl

import (
	"github.com/kube-peering/internal/pkg/connectors"
	"go.uber.org/zap"
)

type Kpctl struct {
	Logger    *zap.SugaredLogger
	VPNConfig connectors.VPNConfig
}

func (ctl *Kpctl) Connect() {
	// TODO pass the ctx and logger, reuse the base logger in sub components
	reverseVPNClient := connectors.NewReverseVPNClient(ctl.VPNConfig)

	reverseVPNClient.Start()
}
