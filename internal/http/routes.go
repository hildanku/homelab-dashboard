package http

import (
	"os/exec"

	"github.com/gofiber/fiber/v2"
	"github.com/hildanku/homelab-dashboard/domain"
	"github.com/hildanku/homelab-dashboard/internal/metrics"
	"github.com/hildanku/homelab-dashboard/internal/services"
	"github.com/hildanku/homelab-dashboard/internal/shared"
	"github.com/hildanku/homelab-dashboard/internal/system"
)

func RegisterRoutes(app *fiber.App, cfg domain.Config) {
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

	app.Get("/api/system", func(c *fiber.Ctx) error {
		sys := system.GetInfo()
		return shared.AppResponse(c, fiber.StatusOK, "success to get system info", sys)
	})

	app.Get("/api/metrics-v2", func(c *fiber.Ctx) error {
		usage, err := metrics.GetUsage()
		if err != nil {
			return shared.AppResponse(c, fiber.StatusInternalServerError, "failed to get usage", nil)
		}
		return shared.AppResponse(c, fiber.StatusOK, "success to get usage", usage)
	})
}
