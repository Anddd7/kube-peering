package kpctl

import (
	"context"
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

/*
req -> backdoorConn -> reqChan -> applicationConn
- 										|
res <- backdoorConn <- resChan <- applicationConn
*/

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

	ctx, cancel := context.WithCancel(context.Background())
	reqChan := make(chan []byte)
	resChan := make(chan []byte)

	go func() {
		err := io.ReadTo(backdoorConn, reqChan)
		logger.Z.Errorf("failed to read from backdoor of peering server: %v", err)
		cancel()
	}()

	go func() {
		err := io.WriteTo(reqChan, applicationConn)
		logger.Z.Errorf("failed to write to target application: %v", err)
		cancel()
	}()

	go func() {
		err := io.ReadTo(applicationConn, resChan)
		logger.Z.Errorf("failed to read from target application: %v", err)
		cancel()
	}()

	go func() {
		err := io.WriteTo(resChan, backdoorConn)
		logger.Z.Errorf("failed to write to backdoor of peering server: %v", err)
		cancel()
	}()

	<-ctx.Done()
}
