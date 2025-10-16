package app

import (
	"log"
	"time"

	"github.com/aronipurwanto/go-download-csv/internal/config"
	httpdeliver "github.com/aronipurwanto/go-download-csv/internal/deliveries/http"
	"github.com/aronipurwanto/go-download-csv/internal/domain/transaction"
	"github.com/aronipurwanto/go-download-csv/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover" // <-- tambahkan ini

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Run() error {
	// DB (Postgres)
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	db, err := gorm.Open(postgres.Open(cfg.DB.DSN()), &gorm.Config{})

	// Auto-migrate
	if err := db.AutoMigrate(&transaction.Transaction{}); err != nil {
		return err
	}

	// Repositories & services
	repo := transaction.NewGormRepository(db)
	service := transaction.NewService(repo)

	// Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "transaction-api",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	app.Use(recover.New()) // OK setelah import recover
	app.Use(cors.New())
	app.Use(middleware.EnforceResponseEnvelope())

	// Router (pakai alias httpdeliver)
	httpdeliver.RegisterRoutes(app, service)

	log.Println("listening on :8080")
	return app.Listen(":8080")
}
