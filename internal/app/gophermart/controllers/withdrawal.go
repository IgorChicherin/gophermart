package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type WithdrawController struct {
}

func (w WithdrawController) Route(api *gin.RouterGroup) {
	withdraw := api.Group("/user/balance")
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
