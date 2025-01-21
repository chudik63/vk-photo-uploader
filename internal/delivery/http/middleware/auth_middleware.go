package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/auth" || c.Request.URL.Path == "/register" || c.Request.URL.Path == "/logout" || strings.Contains(c.Request.URL.Path, "/static") {
			c.Next()
			return
		}

		_, err := c.Cookie("vk_token")

		if err == http.ErrNoCookie {
			c.Redirect(http.StatusUnauthorized, "/auth")
			c.Abort()
			return
		}

		c.Next()
	}
}
