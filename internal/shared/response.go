package shared

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hildanku/homelab-dashboard/domain"
)

func AppResponse(c *fiber.Ctx, status int, message string, result interface{}) error {
	return c.Status(status).JSON(domain.Response{
		Message: message,
		Result:  result,
	})
}
