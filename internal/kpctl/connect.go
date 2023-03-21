package kpctl

import (
	"net"
	"os"

	"github.com/kube-peering/internal/pkg/io"
	"github.com/kube-peering/internal/pkg/logger"
	"github.com/kube-peering/internal/pkg/model"
)

type Kpctl struct {
	Backdoor    model.Backdoor
	Application model.Application
}

func (ctl *Kpctl) Connect() {
	backdoorConn, err := net.Dial("tcp", ctl.Backdoor.Address())
	if err != nil {
		logger.Z.Errorf("failed to connect backdoor of peering server: %v", err)
		os.Exit(1)
	}

	applicationConn, err := net.Dial("tcp", ctl.Application.Address())
	if err != nil {
		logger.Z.Errorf("failed to connect target application: %v", err)
		os.Exit(1)
	}

	io.BiFoward("Backdoor", backdoorConn, "Application", applicationConn)
}
