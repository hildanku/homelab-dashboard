package shared

import "github.com/gofiber/fiber/v2"

type Response struct {
	Message string      `json:"message"`
	Result  interface{} `json:"result,omitempty"`
}

func AppResponse(c *fiber.Ctx, status int, message string, result interface{}) error {
	return c.Status(status).JSON(Response{
		Message: message,
		Result:  result,
	})
}
