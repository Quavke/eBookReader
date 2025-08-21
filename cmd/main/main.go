package main

import (
	"log"

	"github.com/Quavke/eBookReader/pkg/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app, err := config.NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to create app: %s", err)
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}