package main

import (
	"log"
	"ebookr/pkg/config"
)

func main() {
	// Инициализация конфигурации
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Создание и запуск приложения
	app := config.NewApp(cfg)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}