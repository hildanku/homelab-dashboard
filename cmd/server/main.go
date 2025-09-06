package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/hildanku/homelab-dashboard/internal/metrics"
)

func main() {
	app := fiber.New()

	app.Get("/api/metrics", func(c *fiber.Ctx) error {
		snap, err := metrics.SnapshotNow()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(snap)
	})

	log.Println("listening on :5551")
	if err := app.Listen(":5551"); err != nil {
		log.Fatal(err)
	}
}
