package main

import (
	"ebookr/pkg/config"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app := config.NewApp(cfg)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}