package main

import (
	"log"

	"github.com/aronipurwanto/go-download-csv/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
