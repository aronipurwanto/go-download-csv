package http

import (
	"github.com/aronipurwanto/go-download-csv/internal/domain/transaction"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, svc transaction.Service) {
	r := app.Group("/v1")
	RegisterTransactionRoutes(r, svc)
}
