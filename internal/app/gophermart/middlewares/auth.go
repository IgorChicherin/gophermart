package middlewares

import (
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/usecases"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware(userUseCase usecases.UserUseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")

		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"code": "401", "message": "unauthorized"})
			return
		}

		ok, err := userUseCase.Validate(token)

		if err != nil || !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"code": "401", "message": "unauthorized"})
			return
		}

		ctx.Next()
	}
}
