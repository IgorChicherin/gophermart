package controllers

import (
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/middlewares"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WithdrawController struct {
	UserRepository repositories.UserRepository
}

func (w WithdrawController) Route(api *gin.RouterGroup) {
	middleware := middlewares.AuthMiddleware(w.UserRepository)
	withdraw := api.Group("/user/balance").Use(middleware)
	{
		withdraw.GET("/withdrawals", w.withdrawals)
	}
}

// @BasePath /api
// login godoc
// @Summary withdrawals
// @Schemes
// @Description user withdrawals
// @Tags withdrawals
// @Accept json
// @Produce json
// @Success 200 {json} OK
// @Router /user/balance/withdrawals [get]
func (w WithdrawController) withdrawals(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
