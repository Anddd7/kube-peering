package io

import (
	"net"
	"sync"
	"testing"

	"github.com/kube-peering/internal/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestBiForward(t *testing.T) {
	logger.InitSimpleLogger()

	address := ":8080"
	message := "hello, world!"

	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	t.Run("should forward data from source to target", func(t *testing.T) {
		go func() {
			conn1, _ := listener.Accept()
			conn2, _ := listener.Accept()
			BiFoward(conn1, conn2)
		}()

		source, _ := net.Dial("tcp", address)
		defer source.Close()
		target, _ := net.Dial("tcp", address)
		defer target.Close()

		assertWrite(t, source, message)
		assertRead(t, target, message)
	})

	t.Run("should return error if source disconnected", func(t *testing.T) {
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			conn1, _ := listener.Accept()
			conn2, _ := listener.Accept()
			err := BiFoward(conn1, conn2)
			assert.Equal(t, ErrSourceDisconnected, err)
			wg.Done()
		}()

		source, _ := net.Dial("tcp", address)
		target, _ := net.Dial("tcp", address)
		defer target.Close()

		source.Close()

		_, err := source.Write([]byte(message))
		assert.NotNil(t, err)

		wg.Wait()
	})

	t.Run("should return error if target disconnected", func(t *testing.T) {
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			conn1, _ := listener.Accept()
			conn2, _ := listener.Accept()
			err := BiFoward(conn1, conn2)
			assert.Equal(t, ErrTargetDisconnected, err)
			wg.Done()
		}()

		source, _ := net.Dial("tcp", address)
		defer source.Close()
		target, _ := net.Dial("tcp", address)

		target.Close()

		_, err := source.Write([]byte(message))
		assert.Nil(t, err)

		buffer := make([]byte, len(message))
		_, err = target.Read(buffer)
		assert.NotNil(t, err)

		wg.Wait()
	})
}
