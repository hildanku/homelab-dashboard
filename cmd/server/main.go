package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/hildanku/homelab-dashboard/internal/auth"
	"github.com/hildanku/homelab-dashboard/internal/config"
	"github.com/hildanku/homelab-dashboard/internal/http"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())
	cfg := config.Load("config.json")

	authSvc, err := auth.LoadFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	http.RegisterRoutes(app, cfg, authSvc)

	log.Println("listening on :5551")
	if err := app.Listen(":5551"); err != nil {
		log.Fatal(err)
	}
}
