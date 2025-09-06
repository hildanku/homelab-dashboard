package system

import (
	"bufio"
	"os"
	"runtime"
	"strings"

	"github.com/hildanku/homelab-dashboard/domain"
)

func GetInfo() domain.Info {
	info := domain.Info{
		OSName:    "Unknown",
		OSVersion: "Unknown",
		Kernel:    "Unknown",
		Arch:      runtime.GOARCH,
		GoVersion: runtime.Version(),
	}

	// Baca /etc/os-release
	if f, err := os.Open("/etc/os-release"); err == nil {
		defer f.Close()
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			line := sc.Text()
			if strings.HasPrefix(line, "NAME=") {
				info.OSName = strings.Trim(line[5:], `"`)
			}
			if strings.HasPrefix(line, "VERSION=") {
				info.OSVersion = strings.Trim(line[8:], `"`)
			}
		}
	}

	// Baca kernel versi via /proc/sys/kernel/osrelease
	if b, err := os.ReadFile("/proc/sys/kernel/osrelease"); err == nil {
		info.Kernel = strings.TrimSpace(string(b))
	}

	return info
}
