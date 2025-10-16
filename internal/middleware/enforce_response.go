package middleware

import "github.com/gofiber/fiber/v2"

const localEnvelopedKey = "response_enveloped"

func MarkEnveloped(c *fiber.Ctx) { c.Locals(localEnvelopedKey, true) }

func EnforceResponseEnvelope() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.Next(); err != nil {
			return err
		}
		// Skip non-JSON or errors
		ct := c.Get(fiber.HeaderContentType)
		if c.Response().StatusCode() >= 400 || ct == "" || ct == "text/plain" {
			return nil
		}
		if v := c.Locals(localEnvelopedKey); v == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "unenveloped_response"})
		}
		return nil
	}
}
