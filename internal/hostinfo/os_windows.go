package hostinfo

import (
	"fmt"

	"golang.org/x/sys/windows"
)

func osVersion() (name, kernel string) {
	v := windows.RtlGetVersion()
	return fmt.Sprintf("Windows %d.%d build %d", v.MajorVersion, v.MinorVersion, v.BuildNumber), ""
}
