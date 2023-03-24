package io

import (
	"net"

	"github.com/kube-peering/internal/pkg/logger"
)

func StartTCPServer(address string, onReady func(string), onConnected func(net.Conn)) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	onReady(address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Z.Error(err)
			continue
		}
		logger.Z.Infoln("New connection from", conn.RemoteAddr())
		go onConnected(conn)
	}
}
