package http

import (
	"log"
	"net/http"
	"sync"
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
	var wg sync.WaitGroup

	u.jsLog = infrastructure.NewSafeJsonLogger(c)

	path := c.PostForm("path")

	wg.Add(1)

	go func(path string) {
		defer wg.Done()

		if err := u.userService.Send(path); err != nil {
			log.Print(err)
			u.jsLog.SendResponse(http.StatusBadRequest)
			return
		}
		u.jsLog.SendResponse(http.StatusOK)
	}(path)

	wg.Wait()
}
