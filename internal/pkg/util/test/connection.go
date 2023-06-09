package util

import (
	"io"
	"testing"
)

func AssertWrite(t *testing.T, conn io.Writer, data string) {
	t.Helper()
	_, err := conn.Write([]byte(data))
	if err != nil {
		t.Errorf("write failed: %v", err)
	}
}

func AssertRead(t *testing.T, conn io.Reader, expected string) {
	t.Helper()
	buffer := make([]byte, len(expected))
	_, err := conn.Read(buffer)
	if err != nil {
		t.Errorf("read failed: %v", err)
	}

	if string(buffer) != expected {
		t.Errorf("received wrong response: expected %s actual %s", expected, string(buffer))
	}
}
