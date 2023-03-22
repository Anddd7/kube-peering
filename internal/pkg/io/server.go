package io

import (
	"net"

	"github.com/kube-peering/internal/pkg/logger"
)

func AcceptConnections(name, protocol, address string, connChan chan<- net.Conn) {
	ln, err := net.Listen(protocol, address)
	if err != nil {
		logger.Z.Errorf("failed to start %s listener: %v", name, err)
		return
	}
	defer ln.Close()

	for {
		// TODO forbid new connection if exists running connection
		conn, err := ln.Accept()
		if err != nil {
			logger.Z.Error(err)
			continue
		}
		logger.Z.Infof("%s connection is ready", name)
		connChan <- conn
	}
}
