package io

import (
	"context"
	"errors"
	"net"

	"github.com/kube-peering/internal/pkg/logger"
)

var (
	ErrSourceDisconnected = errors.New("source disconnected")
	ErrTargetDisconnected = errors.New("target disconnected")
)

const bufferSize = 1024

func BiFoward(from net.Conn, to net.Conn) error {
	ctx, cancel := context.WithCancel(context.Background())
	errChan := make(chan error)

	go func() {
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				buf := make([]byte, bufferSize)
				_, err := from.Read(buf)
				if err != nil {
					logger.Z.Errorf("=> Got an error %v", err)
					errChan <- ErrSourceDisconnected
					return
				}

				logger.Z.Infof("=> Recive request %s => %s", from.RemoteAddr(), from.LocalAddr())
				logger.Z.Infof("=> Forward to %s => %s", to.LocalAddr(), to.RemoteAddr())

				_, err = to.Write(buf)
				if err != nil {
					logger.Z.Errorf("=> Got an error %v", err)
					errChan <- ErrTargetDisconnected
					return
				}

				logger.Z.Infof("=> %s", string(buf))
			}
		}
	}()

	go func() {
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				buf := make([]byte, bufferSize)
				_, err := to.Read(buf)
				if err != nil {
					logger.Z.Errorf("<= Got an error %v", err)
					errChan <- ErrTargetDisconnected
					return
				}

				logger.Z.Infof("<= Recive response %s => %s", to.RemoteAddr(), to.LocalAddr())
				logger.Z.Infof("<= Forward to %s => %s", from.LocalAddr(), from.RemoteAddr())

				_, err = from.Write(buf)
				if err != nil {
					logger.Z.Errorf("<= Got an error %v", err)
					errChan <- ErrSourceDisconnected
					return
				}

				logger.Z.Infof("<= %s", string(buf))
			}
		}
	}()

	for err := range errChan {
		return err
	}

	return nil
}

func ReadTo(conn net.Conn, buffer chan []byte) error {
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

func WriteTo(buffer chan []byte, conn net.Conn) error {
	for buf := range buffer {
		_, err := conn.Write(buf)
		if err != nil {
			logger.Z.Errorf("Got an error %v", err)
			return err
		}
	}
	return nil
}
