package main

import (
	"errors"

	"github.com/egoist/mygo"
)

type virtualOutput interface {
	write([]float32) error
	close()
}

// StartVirtualMic targets a specific monitor epoch. A delayed setup from a
// canceled session cannot claim the new session's device.
func (m *Morpher) StartVirtualMic(epoch int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if epoch != m.issued*streamSpan || epoch == 0 {
		return errors.New("expired monitor session")
	}
	if m.virtual == nil {
		out, err := openBlackHole()
		if err != nil {
			return err
		}
		m.virtual = out
	}
	return nil
}

func (m *Morpher) StopVirtualMic(epoch int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if epoch == m.issued*streamSpan {
		m.closeVirtualLocked()
	}
}

func (m *Morpher) closeVirtualLocked() {
	if m.virtual != nil {
		m.virtual.close()
		m.virtual = nil
	}
}

func (m *Morpher) BlackHoleGuide() error {
	return mygo.Shell.OpenExternal("https://existential.audio/blackhole/")
}
