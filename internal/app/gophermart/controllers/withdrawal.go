package controllers

import (
	"errors"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/middlewares"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/models"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/repositories"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/usecases"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type WithdrawController struct {
	UserRepository  repositories.UserRepository
	WithdrawUseCase usecases.WithdrawUseCase
}

func (w WithdrawController) Route(api *gin.RouterGroup) {
	middleware := middlewares.AuthMiddleware(w.UserRepository)
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
	token, err := c.Cookie("token")

	if err != nil {
		log.WithFields(log.Fields{"func": "withdrawals"}).Errorln(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized user"})
		return
	}

	login, _, err := w.UserRepository.DecodeToken(token)

	if err != nil {
		log.WithFields(log.Fields{"func": "withdrawals"}).Errorln(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, err := w.UserRepository.GetUser(login)

	if err != nil {
		log.WithFields(log.Fields{"func": "withdrawals"}).Errorln(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	wds, err := w.WithdrawUseCase.WithdrawalsList(user.UserID)

	if err != nil {
		log.WithFields(log.Fields{"func": "withdrawals"}).Errorln(err)
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
	var request models.WithdrawalRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		log.WithFields(log.Fields{"func": "balanceWithdraw"}).Errorln(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	token, err := c.Cookie("token")

	if err != nil {
		log.WithFields(log.Fields{"func": "balanceWithdraw"}).Errorln(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized user"})
		return
	}

	login, _, err := w.UserRepository.DecodeToken(token)

	if err != nil {
		log.WithFields(log.Fields{"func": "balanceWithdraw"}).Errorln(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, err := w.UserRepository.GetUser(login)

	if err != nil {
		log.WithFields(log.Fields{"func": "balanceWithdraw"}).Errorln(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = goluhn.Validate(request.Order)
	if err != nil {
		log.WithFields(log.Fields{"func": "balanceWithdraw"}).Errorln(err)
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	_, err = w.WithdrawUseCase.CreateWithdrawOrder(user, request.Order, request.Sum)

	if errors.Is(err, usecases.ErrInsufficientFunds) {
		log.WithFields(log.Fields{"func": "balanceWithdraw"}).Errorln(err)
		c.AbortWithStatus(http.StatusPaymentRequired)
		return
	}

	if err != nil {
		log.WithFields(log.Fields{"func": "balanceWithdraw"}).Errorln(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
