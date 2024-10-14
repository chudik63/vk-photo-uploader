package http

import (
	"log"
	"net/http"
	"strconv"
	"time"
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

	router.POST("/upload", handler.UploadPhoto)
}

func (p *PhotoHandler) UploadPhoto(c *gin.Context) {
	lastModifiedStr := c.PostForm("lastModified")

	lastModified, _ := strconv.ParseInt(lastModifiedStr, 10, 64)

	modTime := time.Unix(lastModified/1000, 0)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusNoContent, "Ошибка чтения файла")
		log.Fatalf("Ошибка чтения файла: %v", err)
	}
	defer file.Close()

	photo := &entity.Photo{
		File:         file,
		Name:         header.Filename,
		Path:         c.PostForm("path"),
		Size:         header.Size,
		LastModified: modTime,
	}

	if err := p.photoService.UploadPhoto(photo); err != nil {
		log.Printf("Ошибка загрузки файла: %v", err)
		c.String(http.StatusBadRequest, "Ошибка загрузки файла")
		return
	}

	c.String(http.StatusOK, "Фотография загружена")
}
