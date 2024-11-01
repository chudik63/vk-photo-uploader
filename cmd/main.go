package main

import (
	"log"
	"vk-photo-uploader/internal/app"
	"vk-photo-uploader/internal/infrastructure"
)

func main() {
	cfg, err := infrastructure.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	app.Run(cfg)

}
