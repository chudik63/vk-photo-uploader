package http

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"vk-photo-uploader/internal/entity"
	"vk-photo-uploader/internal/service"

	"github.com/gin-gonic/gin"
)

type PhotoHandler struct {
	photoService service.PhotoService
}

func NewPhotoHandler(router *gin.Engine, photoService service.PhotoService) {
	handler := &PhotoHandler{
		photoService: photoService,
	}

	router.POST("/uploader/upload", handler.Upload)
	router.DELETE("/uploader/delete", handler.Delete)
}

func (p *PhotoHandler) Upload(c *gin.Context) {
	count, _ := strconv.Atoi(c.Query("count"))
	folder := c.Query("folder")
	token, _ := c.Cookie("vk_token")

	photos := make([]*entity.Photo, count)

	for i := 0; i < count; i++ {
		file, fileHeader, err := c.Request.FormFile(fmt.Sprintf("file%d", i+1))
		if err != nil {
			c.JSON(http.StatusNoContent, gin.H{"error": "Ошибка чтения файла"})
			log.Fatalf("Ошибка чтения файла: %v", err)
		}
		defer file.Close()

		photos[i] = &entity.Photo{
			File:   file,
			Name:   fileHeader.Filename,
			Folder: folder,
			Size:   fileHeader.Size,
		}
	}

	if err := p.photoService.UploadPhotos(photos, token); err != nil {
		log.Printf("Ошибка загрузки файлов: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка загрузки файлов"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Фотографии загружена"})
}

func (p *PhotoHandler) Delete(c *gin.Context) {
	folderName := c.Query("foldername")
	token, _ := c.Cookie("vk_token")
	if err := p.photoService.DeleteFolder(folderName, token); err != nil {
		log.Printf("Ошибка удаления папки: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка удаления папки"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Папка удалена"})
}
