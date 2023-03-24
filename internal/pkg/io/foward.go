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

				logger.Z.Infof("=> Recive request and forward to target")

				_, err = to.Write(buf)
				if err != nil {
					logger.Z.Errorf("=> Got an error %v", err)
					errChan <- ErrTargetDisconnected
					return
				}
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
				logger.Z.Infof("<= Recive response and forward back to source")

				_, err = from.Write(buf)
				if err != nil {
					logger.Z.Errorf("<= Got an error %v", err)
					errChan <- ErrSourceDisconnected
					return
				}
			}
		}
	}()

	for err := range errChan {
		return err
	}

	return nil
}
