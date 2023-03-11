package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BalanceController struct {
}

func (bc BalanceController) Route(api *gin.RouterGroup) {
	balance := api.Group("/user")
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
