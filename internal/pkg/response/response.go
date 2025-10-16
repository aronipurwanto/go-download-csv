package response

import (
	"github.com/aronipurwanto/go-download-csv/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

type Envelope struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func Success(c *fiber.Ctx, data interface{}, meta interface{}) error {
	middleware.MarkEnveloped(c)
	return c.Status(fiber.StatusOK).JSON(Envelope{Success: true, Message: "OK", Data: data, Meta: meta})
}

func Created(c *fiber.Ctx, data interface{}) error {
	middleware.MarkEnveloped(c)
	return c.Status(fiber.StatusCreated).JSON(Envelope{Success: true, Message: "Created", Data: data})
}

func Error(c *fiber.Ctx, code int, msg string) error {
	middleware.MarkEnveloped(c)
	return c.Status(code).JSON(Envelope{Success: false, Message: msg})
}
