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
	defer backdoorConn.Close()

	applicationConn, err := net.Dial("tcp", ctl.Application.Address())
	if err != nil {
		logger.Z.Errorf("failed to connect target application: %v", err)
		os.Exit(1)
	}
	defer applicationConn.Close()

	err = io.BiFoward(backdoorConn, applicationConn)
	switch err {
	case io.ErrSourceDisconnected:
		logger.Z.Errorf("the connection to backdoor of peering server is closed")
	case io.ErrTargetDisconnected:
		logger.Z.Errorf("the connection to target application is closed")
	}
}
