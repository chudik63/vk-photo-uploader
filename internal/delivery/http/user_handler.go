package http

import (
	"log"
	"net/http"
	"vk-photo-uploader/internal/entity"
	"vk-photo-uploader/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(router *gin.Engine, userService *service.UserService) {
	handler := &UserHandler{
		userService: userService,
	}

	router.POST("/register", handler.Register)
}

func (u *UserHandler) Register(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.String(http.StatusNotAcceptable, "Ошибка чтения данных пользователя")
		log.Fatalf("Ошибка чтения данных пользователя: %v", err)
	}

	u.userService.Register(&user)
}
