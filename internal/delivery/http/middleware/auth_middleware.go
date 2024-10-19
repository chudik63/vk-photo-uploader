package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/auth" || c.Request.URL.Path == "/register" || strings.Contains(c.Request.URL.Path, "/static") {
			c.Next()
			return
		}

		_, err := c.Cookie("vk_token")

		if err == http.ErrNoCookie {
			c.Redirect(http.StatusFound, "/auth")
			c.Abort()
			return
		}

		c.Next()
	}
}
