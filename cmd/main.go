package main

import (
	"log"
	"vk-photo-uploader/internal/delivery/http"
	"vk-photo-uploader/internal/delivery/http/middleware"
	"vk-photo-uploader/internal/infrastructure"
	"vk-photo-uploader/internal/repository"
	"vk-photo-uploader/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := infrastructure.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	router := gin.Default()
	router.Use(middleware.AuthMiddleware())

	vkRepository := repository.NewVkRepository()

	photoService := service.NewPhotoService(vkRepository)

	http.NewPageHandler(router)
	http.NewUserHandler(router)
	http.NewPhotoHandler(router, photoService)

	if err := router.Run(cfg.Server.Port); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
