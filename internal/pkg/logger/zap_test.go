package logger

import "testing"

func TestLogger(t *testing.T) {
	InitLogger(true, "json")

	if Z == nil {
		t.Errorf(`Z should not be nil`)
	}

	Z.Info("log with zap")
}
