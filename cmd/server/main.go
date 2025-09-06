package main

import (
	"log"
	"os/exec"

	"github.com/gofiber/fiber/v2"
	"github.com/hildanku/homelab-dashboard/domain"
	"github.com/hildanku/homelab-dashboard/internal/config"
	"github.com/hildanku/homelab-dashboard/internal/metrics"
	"github.com/hildanku/homelab-dashboard/internal/services"
	"github.com/hildanku/homelab-dashboard/internal/shared"
)

func main() {
	app := fiber.New()

	cfg := config.Load("config.json")

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
			return shared.AppResponse(c, fiber.StatusInternalServerError, "failed to get services", nil)
		}
		return shared.AppResponse(c, fiber.StatusOK, "success to get metrics", stats)
	})

	app.Get("/api/ping/all", func(c *fiber.Ctx) error {
		var out []domain.HTTPStatus
		for _, u := range cfg.HTTPTargets {
			out = append(out, services.PingHTTP(u))
		}
		return shared.AppResponse(c, fiber.StatusOK, "success", out)
	})

	app.Get("/api/docker", func(c *fiber.Ctx) error {
		cmd := "docker ps"
		out, err := exec.Command("sh", "-c", cmd).CombinedOutput()

		return shared.AppResponse(c, fiber.StatusOK, "success to get docker", fiber.Map{
			"ok":     err == nil,
			"output": string(out),
		})
	})

	log.Println("listening on :5551")
	if err := app.Listen(":5551"); err != nil {
		log.Fatal(err)
	}
}
