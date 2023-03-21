package logger

import "testing"

func TestLogger(t *testing.T) {
	if Z != nil {
		t.Errorf(`Z should be nil`)
	}

	InitLogger()

	if Z == nil {
		t.Errorf(`Z should not be nil`)
	}

	Z.Info("log with zap")
}
