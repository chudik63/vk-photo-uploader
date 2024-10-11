package infrastructure

import (
	"sync"

	"github.com/gin-gonic/gin"
)

type SafeJsonLogger struct {
	mu sync.Mutex
	c  *gin.Context
}

func NewSafeJsonLogger(c *gin.Context) *SafeJsonLogger {
	var mu sync.Mutex
	return &SafeJsonLogger{
		mu: mu,
		c:  c,
	}
}

func (l *SafeJsonLogger) SendResponse(status int) {
	l.mu.Lock()
	l.c.JSON(status, gin.H{})
	l.mu.Unlock()
}
