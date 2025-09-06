package services

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/hildanku/homelab-dashboard/domain"
)

func PingHTTP(url string) domain.HTTPStatus {
	client := &http.Client{Timeout: 3 * time.Second}
	start := time.Now()
	resp, err := client.Get(url)
	lat := time.Since(start).Milliseconds()
	if err != nil {
		return domain.HTTPStatus{URL: url, OK: false, Code: 0, Latency: lat}
	}
	defer resp.Body.Close()
	return domain.HTTPStatus{URL: url, OK: resp.StatusCode < 400, Code: resp.StatusCode, Latency: lat}
}

func listComms() (map[string]int, error) {
	counts := map[string]int{}
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if _, err := strconv.Atoi(e.Name()); err != nil {
			continue
		}
		commPath := filepath.Join("/proc", e.Name(), "comm")
		b, err := os.ReadFile(commPath)
		if err != nil {
			continue
		}
		name := strings.TrimSpace(string(b))
		if name != "" {
			counts[name]++
		}
	}
	return counts, nil
}

func CheckProcesses(targets []string) ([]domain.ProcStatus, error) {
	comms, err := listComms()
	if err != nil {
		return nil, err
	}
	out := make([]domain.ProcStatus, 0, len(targets))
	for _, t := range targets {
		c := comms[t]
		out = append(out, domain.ProcStatus{Name: t, Found: c > 0, Count: c})
	}
	return out, nil
}
