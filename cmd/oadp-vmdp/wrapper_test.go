package main

import (
	"testing"

	"github.com/alecthomas/kingpin/v2"
	"github.com/kopia/kopia/cli"
)

func TestNewOADPWrapper(t *testing.T) {
	app := cli.NewApp()
	kp := kingpin.New("test", "test app")

	wrapper := newOADPWrapper(app, kp)

	if wrapper == nil {
		t.Fatal("newOADPWrapper returned nil")
	}

	if wrapper.kopiaApp != app {
		t.Error("wrapper.kopiaApp is not the same as the provided app")
	}

	if wrapper.kp != kp {
		t.Error("wrapper.kp is not the same as the provided kingpin app")
	}
}

func TestOADPWrapperSetup(t *testing.T) {
	app := cli.NewApp()
	kp := kingpin.New("test", "test app")

	wrapper := newOADPWrapper(app, kp)

	// Setup should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("wrapper.setup() panicked: %v", r)
		}
	}()

	wrapper.setup()
}

