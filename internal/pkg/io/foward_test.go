package io

import (
	"bytes"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
