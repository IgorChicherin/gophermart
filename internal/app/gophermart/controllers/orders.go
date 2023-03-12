package controllers

import (
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/middlewares"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OrdersController struct {
	UserRepository repositories.UserRepository
}

func (oc OrdersController) Route(api *gin.RouterGroup) {
	middleware := middlewares.AuthMiddleware(oc.UserRepository)
	orders := api.Group("/user").Use(middleware)
	{
		orders.POST("/orders", oc.orderCreate)
		orders.GET("/orders", oc.orderGet)
	}
}

// @BasePath /api
// login godoc
// @Summary create
// @Schemes
// @Description order create
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {json} OK
// @Router /user/orders [post]
func (oc OrdersController) orderCreate(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

// @BasePath /api
// login godoc
// @Summary get
// @Schemes
// @Description order get
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {json} OK
// @Router /user/orders [get]
func (oc OrdersController) orderGet(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
