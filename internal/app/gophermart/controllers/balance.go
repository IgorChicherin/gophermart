package controllers

import (
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/middlewares"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/repositories"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/usecases"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type BalanceController struct {
	UserRepository repositories.UserRepository
	BalanceUseCase usecases.BalanceUseCase
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
// @Success 200 {json} models.Balance
// @Router /user/balance [get]
func (bc BalanceController) getUserBalance(c *gin.Context) {
	token, err := c.Cookie("token")

	if err != nil {
		log.WithFields(log.Fields{"func": "getUserBalance"}).Errorln(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized user"})
		return
	}

	login, _, err := bc.UserRepository.DecodeToken(token)

	if err != nil {
		log.WithFields(log.Fields{"func": "getUserBalance"}).Errorln(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, err := bc.UserRepository.GetUser(login)

	if err != nil {
		log.WithFields(log.Fields{"func": "getUserBalance"}).Errorln(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	balance, err := bc.BalanceUseCase.GetBalance(user.Login)

	if err != nil {
		log.WithFields(log.Fields{"func": "getUserBalance"}).Errorln(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, balance)
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
