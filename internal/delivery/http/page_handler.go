package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PageHandler struct {
	staticFolder    string
	templatesFolder string
	authPage        string
	uploaderPage    string
}

func NewPageHandler(router *gin.Engine) {
	handler := &PageHandler{
		staticFolder:    "../web/static",
		templatesFolder: "../web/templates",
		authPage:        "auth.html",
		uploaderPage:    "index.html",
	}

	router.LoadHTMLGlob(handler.templatesFolder + "/*.html")
	router.Static("/static", handler.staticFolder)

	router.GET("/", handler.RunUploader)
	router.GET("/auth", handler.RunAuth)
}

func (h *PageHandler) RunAuth(c *gin.Context) {
	c.HTML(http.StatusOK, h.authPage, gin.H{})
}

func (h *PageHandler) RunUploader(c *gin.Context) {
	c.HTML(http.StatusOK, h.uploaderPage, gin.H{})
}
