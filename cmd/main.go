package main

import (
	"log"

	"github.com/IvanKuchsh-600/subtracker/internal/app"
	"github.com/IvanKuchsh-600/subtracker/internal/config"
)

func main() {
	cfg := config.MustLoad()

	app, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	err = app.Run()
	if err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
