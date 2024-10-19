package http

import (
	"log"
	"net/http"
	"vk-photo-uploader/internal/entity"

	"github.com/gin-gonic/gin"
)

type UserHandler struct{}

func NewUserHandler(router *gin.Engine) {
	handler := &UserHandler{}

	router.POST("/register", handler.Register)
	router.POST("/logout", handler.Logout)
}

func (u *UserHandler) Register(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.String(http.StatusNotAcceptable, "Ошибка чтения данных пользователя")
		log.Fatalf("Ошибка чтения данных пользователя: %v", err)
	}
	c.SetCookie("vk_token", user.AccessToken, 1<<20, "/", "localhost", false, true)
}

func (u *UserHandler) Logout(c *gin.Context) {
	c.SetCookie("vk_token", "", -1, "/", "localhost", false, true)
	c.String(http.StatusOK, "Выход")
}
