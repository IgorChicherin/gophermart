package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type OrdersController struct {
}

func (oc OrdersController) Route(api *gin.RouterGroup) {
	orders := api.Group("/user")
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
