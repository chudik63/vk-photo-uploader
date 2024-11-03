package app

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"vk-photo-uploader/internal/delivery/http"
	"vk-photo-uploader/internal/delivery/http/middleware"
	"vk-photo-uploader/internal/infrastructure"
	"vk-photo-uploader/internal/repository"
	"vk-photo-uploader/internal/server"
	"vk-photo-uploader/internal/service"

	"github.com/gin-gonic/gin"
)

func Run(cfg *infrastructure.Config) {
	var wg sync.WaitGroup

	router := gin.Default()

	router.Use(middleware.AuthMiddleware())

	vkRepository := repository.NewVkRepository()

	photoService := service.NewPhotoService(vkRepository)

	http.NewPageHandler(router)
	http.NewUserHandler(router)
	http.NewPhotoHandler(router, photoService, &wg)

	srv := server.NewServer(cfg, router.Handler())

	go func() {
		if err := srv.Run(); err != nil {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	if err := srv.Stop(); err != nil {
		log.Printf("Ошибка остановки сервера: %v", err)
	}

	wg.Wait()

	log.Print("Сервер остановлен")
}
