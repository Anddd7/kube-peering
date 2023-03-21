package kpeering

import (
	"net"

	"github.com/kube-peering/internal/pkg/io"
	"github.com/kube-peering/internal/pkg/logger"
	"github.com/kube-peering/internal/pkg/model"
)

type Kpeering struct {
	Frontdoor model.Frontdoor
	Backdoor  model.Backdoor
}

func (peering *Kpeering) Start() {
	frontdoorListener, err := net.Listen("tcp", peering.Frontdoor.Address())
	if err != nil {
		logger.Z.Errorf("failed to start frontdoor listener: %v", err)
		return
	}
	defer frontdoorListener.Close()

	backdoorListener, err := net.Listen("tcp", peering.Backdoor.Address())
	if err != nil {
		logger.Z.Errorf("failed to start backdoor listener: %v", err)
		return
	}
	defer backdoorListener.Close()

	// TODO use mutex and wait for both front and back conn ready
	for {
		frontdoorConn, err := frontdoorListener.Accept()
		if err != nil {
			logger.Z.Error(err)
			continue
		}
		logger.Z.Infoln("frontdoor connection is comming")

		backdoorConn, err := backdoorListener.Accept()
		if err != nil {
			logger.Z.Error(err)
			continue
		}
		logger.Z.Infoln("backdoor is opened")

		go io.BiFoward("Frontdoor", frontdoorConn, "Backdoor", backdoorConn)
	}
}
