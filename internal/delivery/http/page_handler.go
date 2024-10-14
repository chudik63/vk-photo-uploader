package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PageHandler struct {
	staticFolder    string
	templatesFolder string
	indexPage       string
	uploaderPage    string
}

func NewPageHandler(router *gin.Engine) {
	handler := &PageHandler{
		staticFolder:    "../web/static",
		templatesFolder: "../web/templates",
		indexPage:       "index.html",
		uploaderPage:    "uploader.html",
	}

	router.LoadHTMLGlob(handler.templatesFolder + "/*.html")
	router.Static("/static", handler.staticFolder)

	router.GET("/", handler.RunIndex)
	router.GET("/uploader", handler.RunUploader)
}

func (h *PageHandler) RunIndex(c *gin.Context) {
	c.HTML(http.StatusOK, h.indexPage, gin.H{})
}

func (h *PageHandler) RunUploader(c *gin.Context) {
	c.HTML(http.StatusOK, h.uploaderPage, gin.H{})
}
