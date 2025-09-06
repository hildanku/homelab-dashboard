package metrics

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Snapshot struct {
	CPUUsage  float64 `json:"cpu_usage"`
	RAMTotal  uint64  `json:"ram_total"`
	RAMFree   uint64  `json:"ram_free"`
	RAMUsed   uint64  `json:"ram_used"`
	RAMUsage  float64 `json:"ram_usage"`
	DiskTotal uint64  `json:"disk_total"`
	DiskFree  uint64  `json:"disk_free"`
	DiskUsed  uint64  `json:"disk_used"`
	DiskUsage float64 `json:"disk_usage"`
}

func cpuTotals() (idleAll, total uint64, err error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	if !sc.Scan() {
		return 0, 0, nil
	}

	parts := strings.Fields(sc.Text())
	if len(parts) < 8 || parts[0] != "cpu" {
		return 0, 0, nil
	}

	var user, nice, system, idle, iowait, irq, softirq, steal uint64
	user, _ = strconv.ParseUint(parts[1], 10, 64)
	nice, _ = strconv.ParseUint(parts[2], 10, 64)
	system, _ = strconv.ParseUint(parts[3], 10, 64)
	idle, _ = strconv.ParseUint(parts[4], 10, 64)
	iowait, _ = strconv.ParseUint(parts[5], 10, 64)
	irq, _ = strconv.ParseUint(parts[6], 10, 64)
	softirq, _ = strconv.ParseUint(parts[7], 10, 64)
	if len(parts) > 8 {
		steal, _ = strconv.ParseUint(parts[8], 10, 64)
	}

	idleAll = idle + iowait
	nonIdle := user + nice + system + irq + softirq + steal
	total = idleAll + nonIdle
	return
}

func CPUUsagePercent() (float64, error) {
	idle1, total1, err := cpuTotals()
	if err != nil {
		return 0, nil
	}

	time.Sleep(200 * time.Millisecond)
	idle2, total2, err := cpuTotals()
	if err != nil {
		return 0, err
	}

	totald := float64(total2 - total1)
	idled := float64(idle2 - idle1)
	if totald <= 0 {
		return 0, nil
	}
	return (totald - idled) / totald * 100.0, nil
}

func Memory() (total, used, free uint64, usedPct float64, err error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return
	}
	defer f.Close()

	var memTotal, memFree, buffers, cached, memAvail uint64
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		var key string
		var kb uint64
		_, _ = fmt.Sscanf(sc.Text(), "%s %d kB", &key, &kb)
		b := kb * 1024
		switch key {
		case "MemTotal:":
			memTotal = b
		case "MemFree:":
			memFree = b
		case "Buffers:":
			buffers = b
		case "Cached:":
			cached = b
		case "MemAvailable:":
			memAvail = b
		}
	}

	total = memTotal
	if memAvail > 0 {
		free = memAvail
	} else {
		free = memFree + buffers + cached
	}
	used = total - free
	if total > 0 {
		usedPct = float64(used) / float64(total) * 100.0
	}
	return
}

func DiskUsage(path string) (total, used, free uint64, usedPct float64, err error) {
	var s syscall.Statfs_t
	if err = syscall.Statfs(path, &s); err != nil {
		return
	}
	total = s.Blocks * uint64(s.Bsize)
	free = s.Bfree * uint64(s.Bsize)
	used = total - free
	if total > 0 {
		usedPct = float64(used) / float64(total) * 100.0
	}
	return
}

func SnapshotNow() (Snapshot, error) {
	cpu, err := CPUUsagePercent()
	if err != nil {
		return Snapshot{}, err
	}
	rt, ru, rf, rp, err := Memory()
	if err != nil {
		return Snapshot{}, err
	}
	dt, du, df, dp, err := DiskUsage("/")
	if err != nil {
		return Snapshot{}, err
	}
	return Snapshot{
		CPUUsage:  cpu,
		RAMTotal:  rt,
		RAMUsed:   ru,
		RAMFree:   rf,
		RAMUsage:  rp,
		DiskTotal: dt,
		DiskUsed:  du,
		DiskFree:  df,
		DiskUsage: dp,
	}, nil
}
