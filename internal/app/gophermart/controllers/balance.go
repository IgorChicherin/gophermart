package controllers

import (
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/middlewares"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BalanceController struct {
	UserRepository repositories.UserRepository
}

func (bc BalanceController) Route(api *gin.RouterGroup) {
	middleware := middlewares.AuthMiddleware(bc.UserRepository)
	balance := api.Group("/user").Use(middleware)
	{
		balance.GET("/balance", bc.getUserBalance)
		balance.POST("/balance/withdraw", bc.balanceWithdraw)
	}
}

// @BasePath /api
// login godoc
// @Summary balance
// @Schemes
// @Description get user balance
// @Tags balance
// @Accept json
// @Produce json
// @Success 200 {json} OK
// @Router /user/balance [get]
func (bc BalanceController) getUserBalance(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

// @BasePath /api
// login godoc
// @Summary withdraw
// @Schemes
// @Description withdraw balance
// @Tags balance
// @Accept json
// @Produce json
// @Success 200 {json} OK
// @Router /user/balance/withdraw [post]
func (bc BalanceController) balanceWithdraw(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
