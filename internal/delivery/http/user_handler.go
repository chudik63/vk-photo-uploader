package http

import (
	"log"
	"net/http"
	"vk-photo-uploader/internal/entity"
	"vk-photo-uploader/internal/infrastructure"
	"vk-photo-uploader/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
	jsLog       *infrastructure.SafeJsonLogger
}

func NewUserHandler(router *gin.Engine, userService *service.UserService) {
	handler := &UserHandler{
		userService: userService,
	}

	router.POST("/register", handler.Register)
	router.POST("/send", handler.Send)
}

func (u *UserHandler) Register(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u.userService.Register(&user)
}

func (u *UserHandler) Send(c *gin.Context) {
	path := c.PostForm("path")

	if err := u.userService.Send(path); err != nil {
		log.Print(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
