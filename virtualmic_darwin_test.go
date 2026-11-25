//go:build darwin && !ios

package main

import (
	"os"
	"testing"
)

// The macOS integration job installs the driver on a disposable host. Send
// more than three blocks to exercise completion callbacks and buffer reuse.
func TestBlackHoleDevice(t *testing.T) {
	mode := os.Getenv("VOICEMORPH_TEST_BLACKHOLE")
	if mode == "" {
		t.Skip("set VOICEMORPH_TEST_BLACKHOLE=installed or absent")
	}
	device, err := openBlackHole()
	if mode == "absent" && err != nil {
		t.Log(err)
		return
	}
	if err != nil {
		t.Fatal(err)
	}
	defer device.close()
	if mode != "installed" {
		t.Fatal("unexpected open output", mode)
	}
	for range 16 {
		if err := device.write(make([]float32, blockFrames)); err != nil {
			t.Fatal(err)
		}
	}
}
