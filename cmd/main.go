package main

import (
	"log"
	"vk-photo-uploader/internal/delivery/http"
	"vk-photo-uploader/internal/infrastructure"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := infrastructure.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	router := gin.Default()

	http.NewPageHandler(router)

	if err := router.Run(cfg.Server.Port); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
