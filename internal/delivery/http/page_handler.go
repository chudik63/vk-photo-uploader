package http

import (
	"log"
	"net/http"
	"vk-photo-uploader/internal/entity"

	"github.com/gin-gonic/gin"
)

type PageHandler struct {
	staticFolder    string
	templatesFolder string
	indexPage       string
	folderPage      string
}

func NewPageHandler(router *gin.Engine) {
	handler := &PageHandler{
		staticFolder:    "../web/static",
		templatesFolder: "../web/templates",
		indexPage:       "index.html",
		folderPage:      "folder.html",
	}

	router.LoadHTMLGlob(handler.templatesFolder + "/*.html")
	router.Static("/static", handler.staticFolder)

	router.GET("/", handler.RunIndex)
	router.GET("/folder", handler.RunFolder)

	router.POST("/register", handler.Register)
}

func (h *PageHandler) RunIndex(c *gin.Context) {
	c.HTML(http.StatusOK, h.indexPage, gin.H{})
}

func (h *PageHandler) RunFolder(c *gin.Context) {
	c.HTML(http.StatusOK, h.folderPage, gin.H{})
}

func (u *PageHandler) Register(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
