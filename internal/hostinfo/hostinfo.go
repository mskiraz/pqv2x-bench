// Package hostinfo records the measurement environment of a benchmark run:
// CPU model and features, operating system, Go toolchain and build settings,
// runtime configuration (GOMAXPROCS, garbage collector), and the versions of
// the modules that provide the signature implementations.
//
// Nothing here identifies the machine's owner: no hostname, user name or
// network address is recorded.
package hostinfo

import (
	"os"
	"runtime"
	"runtime/debug"
	"sort"
	"strings"

	"golang.org/x/sys/cpu"
)

// Module is one dependency as recorded in the binary's build information.
type Module struct {
	Path    string `json:"path"`
	Version string `json:"version"`
	Replace string `json:"replace,omitempty"`
}

// Info describes the environment of one run.
type Info struct {
	CPU         string            `json:"cpu"`
	CPUFeatures []string          `json:"cpu_features"`
	NumCPU      int               `json:"num_cpu"`
	GOOS        string            `json:"goos"`
	GOARCH      string            `json:"goarch"`
	OSVersion   string            `json:"os_version"`
	Kernel      string            `json:"kernel,omitempty"`
	GoVersion   string            `json:"go_version"`
	Compiler    string            `json:"compiler"`
	BuildSet    map[string]string `json:"build_settings"`
	GOMAXPROCS  int               `json:"gomaxprocs"`
	GOGCEnv     string            `json:"gogc_env"`
	GCPercent   int               `json:"gc_percent"` // effective; -1 means disabled
	Modules     []Module          `json:"modules"`
}

// GCPercent returns the effective garbage-collector percentage without
// changing it.
func GCPercent() int {
	p := debug.SetGCPercent(100)
	debug.SetGCPercent(p)
	return p
}

// Collect gathers the environment description.
func Collect() Info {
	in := Info{
		CPU:        cpuModel(),
		NumCPU:     runtime.NumCPU(),
		GOOS:       runtime.GOOS,
		GOARCH:     runtime.GOARCH,
		GoVersion:  runtime.Version(),
		Compiler:   runtime.Compiler,
		BuildSet:   map[string]string{},
		GOMAXPROCS: runtime.GOMAXPROCS(0),
		GOGCEnv:    os.Getenv("GOGC"),
		GCPercent:  GCPercent(),
	}
	in.OSVersion, in.Kernel = osVersion()
	in.CPUFeatures = features()
	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, s := range bi.Settings {
			in.BuildSet[s.Key] = s.Value
		}
		for _, d := range bi.Deps {
			m := Module{Path: d.Path, Version: d.Version}
			if d.Replace != nil {
				m.Replace = d.Replace.Path
				if d.Replace.Version != "" {
					m.Replace += "@" + d.Replace.Version
				}
			}
			in.Modules = append(in.Modules, m)
		}
	}
	return in
}

func cpuModel() string {
	if b, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		for _, l := range strings.Split(string(b), "\n") {
			if strings.HasPrefix(l, "model name") || strings.HasPrefix(l, "Model") {
				if i := strings.Index(l, ":"); i >= 0 {
					return strings.TrimSpace(l[i+1:])
				}
			}
		}
		// Many ARM64 kernels print no model name; report the implementer and
		// part numbers instead.
		var impl, part string
		for _, l := range strings.Split(string(b), "\n") {
			if i := strings.Index(l, ":"); i >= 0 {
				k, v := strings.TrimSpace(l[:i]), strings.TrimSpace(l[i+1:])
				if k == "CPU implementer" && impl == "" {
					impl = v
				}
				if k == "CPU part" && part == "" {
					part = v
				}
			}
		}
		if impl != "" || part != "" {
			return "ARM implementer " + impl + ", part " + part
		}
	}
	if v := os.Getenv("PROCESSOR_IDENTIFIER"); v != "" {
		return v
	}
	return "unknown"
}

// features lists the CPU features relevant to the implementations measured
// here (SIMD and hash instructions).
func features() []string {
	var f []string
	add := func(ok bool, name string) {
		if ok {
			f = append(f, name)
		}
	}
	switch runtime.GOARCH {
	case "amd64", "386":
		add(cpu.X86.HasAVX, "avx")
		add(cpu.X86.HasAVX2, "avx2")
		add(cpu.X86.HasAVX512F, "avx512f")
		add(cpu.X86.HasAES, "aes")
		add(cpu.X86.HasBMI2, "bmi2")
	case "arm64":
		add(cpu.ARM64.HasASIMD, "asimd")
		add(cpu.ARM64.HasAES, "aes")
		add(cpu.ARM64.HasSHA2, "sha2")
		add(cpu.ARM64.HasSHA3, "sha3")
		add(cpu.ARM64.HasSHA512, "sha512")
		add(cpu.ARM64.HasATOMICS, "atomics")
	}
	// x/sys/cpu does not expose the x86 SHA extensions; Linux reports them.
	if b, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		for _, l := range strings.Split(string(b), "\n") {
			if strings.HasPrefix(l, "flags") {
				for _, w := range strings.Fields(l) {
					if w == "sha_ni" {
						f = append(f, "sha_ni")
						break
					}
				}
				break
			}
		}
	}
	sort.Strings(f)
	return f
}
