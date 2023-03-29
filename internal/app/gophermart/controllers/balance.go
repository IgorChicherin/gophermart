package controllers

import (
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/usecases"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BalanceController struct {
	UserUseCase    usecases.UserUseCase
	BalanceUseCase usecases.BalanceUseCase
}

func (bc BalanceController) Route(api *gin.RouterGroup, middleware gin.HandlerFunc) {
	balance := api.Group("/user").Use(middleware)
	{
		balance.GET("/balance", bc.getUserBalance)
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
// @Success 200 {json} models.Balance
// @Router /user/balance [get]
func (bc BalanceController) getUserBalance(c *gin.Context) {
	token := c.GetHeader("Authorization")
	user, err := bc.UserUseCase.GetUser(token)

	if err != nil {
		controllerLog(c).WithError(err).Errorln("get user error")
		return
	}

	balance, err := bc.BalanceUseCase.GetBalance(user.Login)

	if err != nil {
		controllerLog(c).WithError(err).Errorln("getting balance error")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, balance)
}
