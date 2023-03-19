package controllers

import (
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/middlewares"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/usecases"
	"github.com/IgorChicherin/gophermart/internal/pkg/accrual"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type OrdersController struct {
	OrderUseCase   usecases.OrderUseCase
	AccrualService accrual.AccrualService
}

func (oc OrdersController) Route(api *gin.RouterGroup) {
	middleware := middlewares.AuthMiddleware(oc.OrderUseCase.GetUserRepository())
	orders := api.Group("/user").Use(middleware)
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
// @Accept plain
// @Produce json
// @Success 200,202
// @Failure 500
// @Failure 400,401,409,422 {object} models.DefaultErrorResponse
// @Router /user/orders [post]
func (oc OrdersController) orderCreate(c *gin.Context) {
	b, err := c.GetRawData()

	if err != nil {
		log.Errorln(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	orderNr := string(b)

	token, err := c.Cookie("token")
	userRepo := oc.OrderUseCase.GetUserRepository()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized user"})
		return
	}

	login, _, err := userRepo.DecodeToken(token)

	if err != nil {
		log.Errorln(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = goluhn.Validate(orderNr)
	if err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	if err != nil {
		log.Errorln(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	orderRepo := oc.OrderUseCase.GetOrderRepository()
	hasOrder, err := orderRepo.HasOrder(orderNr)

	if err != nil {
		log.Errorln(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if hasOrder {
		order, err := orderRepo.GetOrder(orderNr)

		if err != nil {
			log.Errorln(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		user, err := userRepo.GetUser(login)

		if err != nil {
			log.Errorln(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		if order.UserID != user.UserID {
			c.AbortWithStatus(http.StatusConflict)
			return
		}

		c.AbortWithStatus(http.StatusOK)
		return
	}

	_, err = oc.OrderUseCase.CreateOrder(login, orderNr)

	if err != nil {
		log.Errorln(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
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
	token, err := c.Cookie("token")
	userRepo := oc.OrderUseCase.GetUserRepository()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized user"})
		return
	}

	login, _, err := userRepo.DecodeToken(token)

	if err != nil {
		log.Errorln(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ordersList, err := oc.OrderUseCase.GetOrdersList(login)
	if err != nil {
		log.Errorln(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, ordersList)
}
