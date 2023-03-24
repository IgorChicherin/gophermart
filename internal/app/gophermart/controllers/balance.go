package controllers

import (
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/middlewares"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/repositories"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/usecases"
	"github.com/gin-gonic/gin"
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

	if token == "" {
		controllerLog(c).Errorln("unauthorized")
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized user"})
		return
	}

	login, _, err := bc.UserRepository.DecodeToken(token)

	if err != nil {
		controllerLog(c).WithError(err).Errorln("can't decode token")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, err := bc.UserRepository.GetUser(login)

	if err != nil {
		controllerLog(c).WithError(err).Errorln("getting user error")
		c.AbortWithStatus(http.StatusInternalServerError)
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
