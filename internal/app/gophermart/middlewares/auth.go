package middlewares

import (
	"errors"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/repositories"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func AuthMiddleware(userRepo repositories.UserRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("token")

		if err != nil && !errors.Is(err, http.ErrNoCookie) {
			log.Errorf("auth middleware error: %s", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

		if err != nil && errors.Is(err, http.ErrNoCookie) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"code": "401", "message": "unauthorized"})
			return
		}

		ok, err := userRepo.Validate(token)

		if err != nil || !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"code": "401", "message": "unauthorized"})
			return
		}

		ctx.Next()
	}
}
