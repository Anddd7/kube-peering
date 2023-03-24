package io

import (
	"net"
	"sync"
	"testing"

	util_test "github.com/kube-peering/internal/pkg/util/test"
	"github.com/stretchr/testify/assert"
)

func TestStartTCPServer(t *testing.T) {
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
			util_test.AssertRead(t, conn, request)
			util_test.AssertWrite(t, conn, response)
		},
	)

	wg.Wait()

	conn, err := net.Dial("tcp", address)
	if err != nil {
		t.Errorf("dial failed: %v", err)
	}
	defer conn.Close()

	util_test.AssertWrite(t, conn, request)
	util_test.AssertRead(t, conn, response)
}
