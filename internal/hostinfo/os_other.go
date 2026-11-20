//go:build !linux && !windows

package hostinfo

import "runtime"

func osVersion() (name, kernel string) { return runtime.GOOS, "" }
