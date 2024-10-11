package http

import (
	"log"
	"net/http"
	"strconv"
	"time"
	"vk-photo-uploader/internal/entity"
	"vk-photo-uploader/internal/infrastructure"
	"vk-photo-uploader/internal/service"

	"github.com/gin-gonic/gin"
)

type PhotoHandler struct {
	photoService service.PhotoService
	jsLog        *infrastructure.SafeJsonLogger
}

func NewPhotoHandler(router *gin.Engine, photoService service.PhotoService) {
	handler := &PhotoHandler{
		photoService: photoService,
	}

	router.POST("/upload", handler.UploadPhoto)
}

func (p *PhotoHandler) UploadPhoto(c *gin.Context) {
	lastModifiedStr := c.PostForm("lastModified")

	lastModified, err := strconv.ParseInt(lastModifiedStr, 10, 64)
	if err != nil {
		log.Print(err)
	}
	modTime := time.Unix(lastModified/1000, 0)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Error retrieving file")
		return
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
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
