package http

import (
	"fmt"
	"log"
	"net/http"
	"vk-photo-storage/internal/entity"

	"github.com/gin-gonic/gin"
)

type PageHandler struct {
	root       string
	indexPage  string
	folderPage string
}

func NewPageHandler(router *gin.Engine) {
	handler := &PageHandler{
		root:       "../web/static/",
		indexPage:  "index.html",
		folderPage: "folder.html",
	}

	handler.SetRoot(router)

	router.GET("/", handler.RunIndex)
	router.GET("/folder", handler.RunFolder)

	router.POST("/register", handler.Register)
	router.POST("/select", handler.Select)
}

func (h *PageHandler) SetRoot(r *gin.Engine) {
	r.LoadHTMLGlob(fmt.Sprintf("%s*", h.root))
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

func (u *PageHandler) Select(c *gin.Context) {
	var path entity.Path
	if err := c.ShouldBindJSON(&path); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Print(path.Path)
}
