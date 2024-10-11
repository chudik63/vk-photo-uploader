package http

import (
	"log"
	"net/http"
	"strconv"
	"sync"
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
	var wg sync.WaitGroup

	p.jsLog = infrastructure.NewSafeJsonLogger(c)

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

	wg.Add(1)

	go func(photo *entity.Photo) {
		defer wg.Done()

		if err := p.photoService.UploadPhoto(photo); err != nil {
			log.Print(err)
			return
		}
		p.jsLog.SendResponse(http.StatusOK)
		log.Print(photo.Name, " CREATED")
	}(photo)

	wg.Wait()
}
