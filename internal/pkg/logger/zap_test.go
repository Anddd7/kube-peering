package logger

import "testing"

func TestLogger(t *testing.T) {
	if Z != nil {
		t.Errorf(`Z should be nil`)
	}

	InitSimpleLogger()

	if Z == nil {
		t.Errorf(`Z should not be nil`)
	}

	Z.Info("log with zap")
}
