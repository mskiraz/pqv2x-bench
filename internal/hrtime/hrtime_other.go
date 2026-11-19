//go:build !windows

// Package hrtime provides a monotonic nanosecond clock with sub-microsecond
// resolution.
package hrtime

import "time"

var origin = time.Now()

// Now returns nanoseconds from an arbitrary origin.
func Now() int64 { return int64(time.Since(origin)) }

// Resolution returns the nominal clock tick in nanoseconds.
func Resolution() int64 { return 1 }
