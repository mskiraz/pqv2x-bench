//go:build windows

// Package hrtime provides a monotonic nanosecond clock with sub-microsecond
// resolution. On Windows, time.Now advances in steps of about 0.5 ms, which
// is too coarse for operations of tens of microseconds, so the performance
// counter is read directly.
package hrtime

import (
	"syscall"
	"unsafe"
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	procQPC  = kernel32.NewProc("QueryPerformanceCounter")
	procQPF  = kernel32.NewProc("QueryPerformanceFrequency")
	freq     int64
)

func init() {
	procQPF.Call(uintptr(unsafe.Pointer(&freq)))
	if freq <= 0 {
		panic("hrtime: QueryPerformanceFrequency failed")
	}
}

// Now returns nanoseconds from an arbitrary origin.
func Now() int64 {
	var c int64
	syscall.SyscallN(procQPC.Addr(), uintptr(unsafe.Pointer(&c)))
	return (c/freq)*1_000_000_000 + (c%freq)*1_000_000_000/freq
}

// Resolution returns the clock tick in nanoseconds.
func Resolution() int64 { return 1_000_000_000 / freq }
