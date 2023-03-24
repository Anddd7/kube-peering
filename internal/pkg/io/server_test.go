package io

import (
	"net"
	"sync"
	"testing"

	"github.com/kube-peering/internal/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestStartTCPServer(t *testing.T) {
	logger.InitSimpleLogger()

	wg := sync.WaitGroup{}
	wg.Add(1)

	address := ":8080"
	request := "hello, world!"
	response := "hello, world, too!"

	go StartTCPServer(address,
		func(addr string) {
			assert.Equal(t, address, addr)
			wg.Done()
		},
		func(conn net.Conn) {
			assertRead(t, conn, request)
			assertWrite(t, conn, response)
		},
	)

	wg.Wait()

	conn, err := net.Dial("tcp", address)
	if err != nil {
		t.Errorf("dial failed: %v", err)
	}
	defer conn.Close()

	assertWrite(t, conn, request)
	assertRead(t, conn, response)
}

func assertWrite(t *testing.T, conn net.Conn, data string) {
	_, err := conn.Write([]byte(data))
	if err != nil {
		t.Errorf("write failed: %v", err)
	}
}

func assertRead(t *testing.T, conn net.Conn, expected string) {
	buffer := make([]byte, len(expected))
	_, err := conn.Read(buffer)
	if err != nil {
		t.Errorf("read failed: %v", err)
	}

	if string(buffer) != expected {
		t.Errorf("received wrong response: expected %s actual %s", expected, string(buffer))
	}
}
