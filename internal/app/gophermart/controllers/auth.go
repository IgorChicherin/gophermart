package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController struct {
}

func (ac AuthController) Route(api *gin.RouterGroup) {
	auth := api.Group("/user")
	{
		auth.POST("/login", ac.login)
		auth.POST("/register", ac.register)
	}
}

// @BasePath /api
// login godoc
// @Summary login
// @Schemes
// @Description user login
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {json} OK
// @Router /auth/login [post]
func (ac AuthController) login(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

// @BasePath /api
// login godoc
// @Summary register
// @Schemes
// @Description user registration
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {json} OK
// @Router /auth/register [post]
func (ac AuthController) register(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
