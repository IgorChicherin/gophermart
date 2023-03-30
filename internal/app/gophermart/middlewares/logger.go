package middlewares

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("logger", log.WithField("method", c.FullPath()))
		c.Next()
	}
}
