package io

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	"testing"

	util_test "github.com/kube-peering/internal/pkg/util/test"
	"github.com/stretchr/testify/assert"
)

func TestBiForward(t *testing.T) {
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

		util_test.AssertWrite(t, source, message)
		util_test.AssertRead(t, target, message)
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

func TestReadTo(t *testing.T) {
	testString := "hello world"
	expectedBytes := bytes.Buffer{}
	expectedBytes.WriteString(testString)

	// setup mock server
	listener, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		conn, _ := listener.Accept()
		conn.Write(expectedBytes.Bytes())
		conn.Close()
	}()

	fmt.Println(listener.Addr().String())

	// setup mock connection
	buffer := make(chan []byte)
	conn, _ := net.Dial("tcp", listener.Addr().String())

	// when
	go func() {
		err := ReadTo(conn, buffer)

		assert.NotNil(t, err, "read loop has ended when connection closed")
	}()

	// post verify
	actualBytes := bytes.Buffer{}
	actualBytes.Write(<-buffer)
	if !bytes.Equal(expectedBytes.Bytes(), actualBytes.Bytes()) {
		t.Errorf("Expected %v but got %v", expectedBytes.Bytes(), actualBytes.Bytes())
	}
}

func TestWriteTo(t *testing.T) {
	testString := "hello world"
	expectedBytes := bytes.Buffer{}
	expectedBytes.WriteString(testString)

	receivedBuffer := make(chan []byte)

	// setup mock server
	listener, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		conn, _ := listener.Accept()
		receivedBytes := make([]byte, len(expectedBytes.Bytes()))
		conn.Read(receivedBytes)
		conn.Close()
		receivedBuffer <- receivedBytes
	}()

	// setup mock connection
	buffer := make(chan []byte)
	conn, _ := net.Dial("tcp", listener.Addr().String())

	// when
	go func() {
		err := WriteTo(buffer, conn)

		assert.NotNil(t, err, "write loop has ended when connection closed")
	}()

	// post verify
	buffer <- expectedBytes.Bytes()
	receivedBytes := bytes.Buffer{}
	receivedBytes.Write(<-receivedBuffer)
	if !bytes.Equal(expectedBytes.Bytes(), receivedBytes.Bytes()) {
		t.Errorf("Expected %v but got %v", expectedBytes.Bytes(), receivedBytes.Bytes())
	}
}
