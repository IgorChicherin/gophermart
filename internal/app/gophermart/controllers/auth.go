package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController struct {
}

func (ac AuthController) Route(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.GET("/login/", ac.Login)
	}
}

// @BasePath /api
// Login godoc
// @Summary login
// @Schemes
// @Description do ping
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {string} "OK"
// @Router /auth/login/ [get]
func (ac AuthController) Login(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
