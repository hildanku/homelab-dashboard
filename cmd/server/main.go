package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/hildanku/homelab-dashboard/domain"
	"github.com/hildanku/homelab-dashboard/internal/metrics"
	"github.com/hildanku/homelab-dashboard/internal/services"
	"github.com/hildanku/homelab-dashboard/internal/shared"
)

func loadConfig(path string) domain.Config {
	var cfg domain.Config
	b, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}
	_ = json.Unmarshal(b, &cfg)
	return cfg
}

func main() {
	app := fiber.New()

	cfg := loadConfig("config.json")

	app.Get("/api/metrics", func(c *fiber.Ctx) error {
		snap, err := metrics.SnapshotNow()
		if err != nil {
			return shared.AppResponse(c, fiber.StatusInternalServerError, "failed to get metrics", nil)
		}
		return shared.AppResponse(c, fiber.StatusOK, "success to get metrics", snap)
	})

	app.Get("/api/services", func(c *fiber.Ctx) error {
		targets := cfg.ProcessTargets
		if len(targets) == 0 {
			targets = []string{"dockerd"} // fallback
		}
		stats, err := services.CheckProcesses(targets)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(stats)
	})

	app.Get("/api/ping/all", func(c *fiber.Ctx) error {
		var out []services.HTTPStatus
		for _, u := range cfg.HTTPTargets {
			out = append(out, services.PingHTTP(u))
		}
		return c.JSON(out)
	})

	log.Println("listening on :5551")
	if err := app.Listen(":5551"); err != nil {
		log.Fatal(err)
	}
}
