package io

import (
	"net"

	"github.com/kube-peering/internal/pkg/logger"
)

func BiFoward(fromName string, from net.Conn, toName string, to net.Conn) {
	defer from.Close()
	defer to.Close()

	request := make(chan []byte)
	response := make(chan []byte)
	bufferSize := 1024

	go func() {
		for {
			buf := make([]byte, bufferSize)
			n, err := from.Read(buf)
			if err != nil {
				logger.Z.Errorf("=> Got an error %v", err)
				break
			}
			logger.Z.Infof("=> Recive request data from %s", fromName)
			request <- buf[:n]
		}
	}()

	go func() {
		for {
			buf := make([]byte, bufferSize)
			n, err := to.Read(buf)
			if err != nil {
				logger.Z.Errorf("<= Got an error %v", err)
				break
			}
			logger.Z.Infof("<= Recive response data from %s", toName)
			response <- buf[:n]
		}
	}()

	for {
		select {
		case buf := <-request:
			_, err := to.Write(buf)
			logger.Z.Infof("=> Forward request to %s", toName)
			if err != nil {
				logger.Z.Errorf("=> Got an error %v", err)
			}
		case buf := <-response:
			_, err := from.Write(buf)
			logger.Z.Infof("<= Forward response to %s", fromName)
			if err != nil {
				logger.Z.Errorf("<= Got an error %v", err)
			}
		}
	}
}
