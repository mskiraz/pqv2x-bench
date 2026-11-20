package hostinfo

import (
	"os"
	"strings"
)

func osVersion() (name, kernel string) {
	if b, err := os.ReadFile("/etc/os-release"); err == nil {
		for _, l := range strings.Split(string(b), "\n") {
			if v, ok := strings.CutPrefix(l, "PRETTY_NAME="); ok {
				name = strings.Trim(v, `"`)
			}
		}
	}
	if b, err := os.ReadFile("/proc/sys/kernel/osrelease"); err == nil {
		kernel = strings.TrimSpace(string(b))
	}
	if name == "" {
		name = "linux"
	}
	return name, kernel
}
