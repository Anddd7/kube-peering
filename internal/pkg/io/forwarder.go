package io

import (
	"context"
	"net"
	"os"

	"github.com/kube-peering/internal/pkg/logger"
	"github.com/kube-peering/internal/pkg/model"
)

/*
forwardChan  -> TCPForwarder -> server (conn)
backwordChan <- TCPForwarder <- server (conn)
*/
type TCPForwarder struct {
	ctx          context.Context
	address      string
	forwardChan  chan []byte
	backwordChan chan []byte
}

func NewTCPForwarder(
	ctx context.Context,
	cfg model.Forwarder,
	forwardChan chan []byte,
	backwordChan chan []byte,
) *TCPForwarder {
	return &TCPForwarder{
		ctx:          context.WithValue(ctx, keyComponentID, cfg.Name),
		address:      cfg.Address(),
		forwardChan:  forwardChan,
		backwordChan: backwordChan,
	}
}

func (t *TCPForwarder) name() string {
	return t.ctx.Value("component").(string)
}

func (t *TCPForwarder) Run() {
	_, cancel := context.WithCancel(t.ctx)

	conn, err := net.Dial("tcp", t.address)
	if err != nil {
		logger.Z.Errorf("[%s] failed to connect the server: %v", t.name(), err)
		os.Exit(1)
	}
	defer conn.Close()

	go func() {
		err := WriteTo(t.forwardChan, conn)
		logger.Z.Errorf("[%s] failed to write to server: %v", t.name(), err)
		cancel()
	}()

	go func() {
		err := ReadTo(conn, t.backwordChan)
		logger.Z.Errorf("[%s] failed to read from server: %v", t.name(), err)
		cancel()
	}()

	<-t.ctx.Done()
}
