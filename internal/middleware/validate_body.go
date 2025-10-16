package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func ValidateBody[T any](validate func(T) error, localKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req T
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid_json", "detail": err.Error()})
		}
		if validate != nil {
			if err := validate(req); err != nil {
				return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "validation_failed", "detail": err.Error()})
			}
		}
		c.Locals(localKey, req)
		return c.Next()
	}
}
