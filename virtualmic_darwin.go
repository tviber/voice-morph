//go:build darwin && !ios

package main

import (
	"fmt"
	"sync"
	"time"
	"unsafe"

	"github.com/ebitengine/purego"
)

// AudioToolbox's AudioQueue tracks a specific device when CurrentDevice is
// set to its UID. No process-wide or system-wide audio setting is changed.
// Both supported macOS architectures use the same 64-bit Core Audio ABI.
type streamFormat struct {
	rate                                                                           float64
	format, flags, packetBytes, packetFrames, frameBytes, channels, bits, reserved uint32
}
type queueBuffer struct {
	capacity uint32
	data     unsafe.Pointer
	size     uint32
}

var aq struct {
	deviceForUID    func(uint32, *[3]uint32, uint32, *uintptr, *uint32, *uint32) int32
	newOutput       func(*streamFormat, uintptr, uintptr, uintptr, uintptr, uint32, *uintptr) int32
	setProperty     func(uintptr, uint32, *uintptr, uint32) int32
	allocate        func(uintptr, uint32, **queueBuffer) int32
	enqueue         func(uintptr, *queueBuffer, uint32, uintptr) int32
	start           func(uintptr, uintptr) int32
	dispose         func(uintptr, bool) int32
	newString       func(uintptr, string, uint32) uintptr
	release         func(uintptr)
	semaphore       func(int) uintptr
	wait            func(uintptr, uint64) int
	deadline        func(uint64, int64) uint64
	dispatchRelease func(uintptr)
	signal          uintptr
}

var loadQueue = sync.OnceValue(func() error {
	lib, err := purego.Dlopen("/System/Library/Frameworks/AudioToolbox.framework/AudioToolbox", purego.RTLD_LAZY|purego.RTLD_GLOBAL)
	if err != nil {
		return err
	}
	purego.RegisterLibFunc(&aq.newOutput, lib, "AudioQueueNewOutput")
	purego.RegisterLibFunc(&aq.setProperty, lib, "AudioQueueSetProperty")
	purego.RegisterLibFunc(&aq.deviceForUID, lib, "AudioObjectGetPropertyData")
	purego.RegisterLibFunc(&aq.allocate, lib, "AudioQueueAllocateBuffer")
	purego.RegisterLibFunc(&aq.enqueue, lib, "AudioQueueEnqueueBuffer")
	purego.RegisterLibFunc(&aq.start, lib, "AudioQueueStart")
	purego.RegisterLibFunc(&aq.dispose, lib, "AudioQueueDispose")
	purego.RegisterLibFunc(&aq.newString, lib, "CFStringCreateWithCString")
	purego.RegisterLibFunc(&aq.release, lib, "CFRelease")
	purego.RegisterLibFunc(&aq.semaphore, purego.RTLD_DEFAULT, "dispatch_semaphore_create")
	purego.RegisterLibFunc(&aq.wait, purego.RTLD_DEFAULT, "dispatch_semaphore_wait")
	purego.RegisterLibFunc(&aq.deadline, purego.RTLD_DEFAULT, "dispatch_time")
	purego.RegisterLibFunc(&aq.dispatchRelease, purego.RTLD_DEFAULT, "dispatch_release")
	// The queue passes userData (our semaphore) first. As in MyGo's audio
	// backend, dispatch_semaphore_signal ignores the remaining callback args.
	aq.signal, err = purego.Dlsym(purego.RTLD_DEFAULT, "dispatch_semaphore_signal")
	return err
})

const blockFrames = 2048 // public/worklet.js capture block

type blackHoleQueue struct {
	queue, semaphore uintptr
	buffers          [3]*queueBuffer
	next             int
	started          bool
}

func openBlackHole() (virtualOutput, error) {
	if err := loadQueue(); err != nil {
		return nil, err
	}
	uid := aq.newString(0, "BlackHole2ch_UID", 0x08000100) // UTF-8
	defer aq.release(uid)
	address := [3]uint32{0x75696464, 0x676c6f62, 0} // TranslateUIDToDevice, global, main
	var deviceID uint32
	size := uint32(unsafe.Sizeof(deviceID))
	if status := aq.deviceForUID(1, &address, uint32(unsafe.Sizeof(uid)), &uid, &size, &deviceID); status != 0 || deviceID == 0 {
		return nil, fmt.Errorf("BlackHole 2ch unavailable (%d)", status)
	}
	d := &blackHoleQueue{semaphore: aq.semaphore(0)}
	format := streamFormat{rate: 48000, format: 0x6c70636d, flags: 1 | 8, packetBytes: 8, packetFrames: 1, frameBytes: 8, channels: 2, bits: 32} // packed stereo float LPCM
	if status := aq.newOutput(&format, aq.signal, d.semaphore, 0, 0, 0, &d.queue); status != 0 {
		d.close()
		return nil, fmt.Errorf("open audio queue (%d)", status)
	}
	if status := aq.setProperty(d.queue, 0x61716364, &uid, uint32(unsafe.Sizeof(uid))); status != 0 { // 'aqcd'
		d.close()
		return nil, fmt.Errorf("BlackHole 2ch unavailable (%d)", status)
	}
	for i := range d.buffers {
		if status := aq.allocate(d.queue, blockFrames*8, &d.buffers[i]); status != 0 {
			d.close()
			return nil, fmt.Errorf("allocate output buffer (%d)", status)
		}
	}
	return d, nil
}

func (d *blackHoleQueue) write(samples []float32) error {
	for len(samples) > 0 {
		if d.started && aq.wait(d.semaphore, aq.deadline(0, int64(300*time.Millisecond))) != 0 {
			return fmt.Errorf("BlackHole 2ch stopped consuming audio")
		}
		n := min(blockFrames, len(samples))
		buf := d.buffers[d.next]
		stereo := unsafe.Slice((*float32)(buf.data), 2*n)
		for i, sample := range samples[:n] {
			stereo[2*i], stereo[2*i+1] = sample, sample
		}
		buf.size = uint32(n * 8)
		if status := aq.enqueue(d.queue, buf, 0, 0); status != 0 {
			return fmt.Errorf("enqueue audio (%d)", status)
		}
		d.next = (d.next + 1) % len(d.buffers)
		if !d.started && d.next == 0 {
			if status := aq.start(d.queue, 0); status != 0 {
				return fmt.Errorf("start BlackHole 2ch (%d)", status)
			}
			d.started = true
		}
		samples = samples[n:]
	}
	return nil
}

func (d *blackHoleQueue) close() {
	if d.queue != 0 {
		aq.dispose(d.queue, true)
	}
	aq.dispatchRelease(d.semaphore)
}
