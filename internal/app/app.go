package app

import (
	"log"
	"os"
	"os/signal"
	"vk-photo-uploader/internal/delivery/http"
	"vk-photo-uploader/internal/delivery/http/middleware"
	"vk-photo-uploader/internal/infrastructure"
	"vk-photo-uploader/internal/repository"
	"vk-photo-uploader/internal/server"
	"vk-photo-uploader/internal/service"

	"github.com/gin-gonic/gin"
)

func Run(cfg *infrastructure.Config) {
	router := gin.Default()

	router.Use(middleware.AuthMiddleware())

	vkRepository := repository.NewVkRepository()

	photoService := service.NewPhotoService(vkRepository)

	http.NewPageHandler(router)
	http.NewUserHandler(router)
	http.NewPhotoHandler(router, photoService)

	srv := server.NewServer(cfg, router.Handler())

	go func() {
		if err := srv.Run(); err != nil {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	if err := srv.Stop(); err != nil {
		log.Printf("Ошибка остановки сервера: %v", err)
	}
}
