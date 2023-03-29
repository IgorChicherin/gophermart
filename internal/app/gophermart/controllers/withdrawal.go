package controllers

import (
	"errors"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/models"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/usecases"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WithdrawController struct {
	UserUseCase     usecases.UserUseCase
	WithdrawUseCase usecases.WithdrawUseCase
}

func (w WithdrawController) Route(api *gin.RouterGroup, middleware gin.HandlerFunc) {
	withdraw := api.Group("/user").Use(middleware)
	{
		withdraw.GET("/withdrawals", w.withdrawals)
		withdraw.POST("/balance/withdraw", w.balanceWithdraw)
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
// @Success 200 {json} []models.Withdraw
// @Success 204
// @Router /user/withdrawals [get]
func (w WithdrawController) withdrawals(c *gin.Context) {
	token := c.GetHeader("Authorization")
	user, err := w.UserUseCase.GetUser(token)

	if err != nil {
		controllerLog(c).WithError(err).Errorln("get user error")
		return
	}

	wds, err := w.WithdrawUseCase.WithdrawalsList(user.UserID)

	if err != nil {
		controllerLog(c).WithError(err).Errorln("withdrawals list error")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if len(wds) == 0 {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, wds)
}

// @BasePath /api
// login godoc
// @Summary withdraw
// @Schemes
// @Description withdraw balance
// @Tags balance
// @Tags withdrawals
// @Accept json
// @Produce json
// @Success 200 {json} OK
// @Router /user/balance/withdraw [post]
func (w WithdrawController) balanceWithdraw(c *gin.Context) {
	token := c.GetHeader("Authorization")
	user, err := w.UserUseCase.GetUser(token)

	if err != nil {
		controllerLog(c).WithError(err).Errorln("get user error")
		return
	}

	var request models.WithdrawalRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		controllerLog(c).WithError(err).Errorln("can't parse request")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = goluhn.Validate(request.Order)
	if err != nil {
		controllerLog(c).WithError(err).Errorln("order is not valid")
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	_, err = w.WithdrawUseCase.CreateWithdrawOrder(user, request.Order, request.Sum)

	if errors.Is(err, usecases.ErrInsufficientFunds) {
		controllerLog(c).WithError(err)
		c.AbortWithStatus(http.StatusPaymentRequired)
		return
	}

	if err != nil {
		controllerLog(c).WithError(err).Errorln("couldn't withdraw")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
