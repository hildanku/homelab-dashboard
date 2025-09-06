package config

import (
	"encoding/json"
	"os"

	"github.com/hildanku/homelab-dashboard/domain"
)

func Load(path string) domain.Config {
	var cfg domain.Config
	b, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}
	_ = json.Unmarshal(b, &cfg)
	return cfg
}
