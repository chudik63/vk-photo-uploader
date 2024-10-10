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
	_ = handler

	router.POST("/upload", handler.Upload)
}

func (p *PhotoHandler) Upload(c *gin.Context) {
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

	size, err := strconv.ParseUint(c.PostForm("size"), 10, 64)
	if err != nil {
		log.Print(err)
	}

	photo := &entity.Photo{
		File:         file,
		Name:         c.PostForm("name"),
		Path:         c.PostForm("path"),
		Size:         size,
		Type:         c.PostForm("type"),
		LastModified: modTime,
	}

	_, _ = photo, header

	// savePath := "./uploads/"
	// err = os.MkdirAll(savePath, os.ModePerm)
	// if err != nil {
	// 	c.String(http.StatusInternalServerError, "Error creating directory")
	// 	return
	// }

	// dst, err := os.Create(filepath.Join(savePath, header.Filename))
	// if err != nil {
	// 	c.String(http.StatusInternalServerError, "Error creating file")
	// 	return
	// }
	// defer dst.Close()

	// if _, err := io.Copy(dst, file); err != nil {
	// 	c.String(http.StatusInternalServerError, "Error saving file")
	// 	return
	// }

	// if ok c.JSON(http.StatusOK, gin.H{})

}
