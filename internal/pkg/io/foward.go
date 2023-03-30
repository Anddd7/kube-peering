package io

import (
	"io"

	"github.com/kube-peering/internal/pkg/logger"
)

const bufferSize = 1024

func ReadTo(conn io.Reader, buffer chan []byte) error {
	for {
		buf := make([]byte, bufferSize)
		n, err := conn.Read(buf)
		if err != nil {
			logger.Z.Errorf("Got an error %v", err)
			return err
		}
		buffer <- buf[:n]
	}
}

func WriteTo(buffer chan []byte, conn io.Writer) error {
	for buf := range buffer {
		_, err := conn.Write(buf)
		if err != nil {
			logger.Z.Errorf("Got an error %v", err)
			return err
		}
	}
	return nil
}
