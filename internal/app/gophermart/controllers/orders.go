package controllers

import (
	"net/http"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/gin-gonic/gin"

	"github.com/IgorChicherin/gophermart/internal/app/gophermart/middlewares"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/usecases"
	"github.com/IgorChicherin/gophermart/internal/pkg/accrual"
)

type OrdersController struct {
	UserUseCase    usecases.UserUseCase
	OrderUseCase   usecases.OrderUseCase
	AccrualService accrual.AccrualService
}

func (oc OrdersController) Route(api *gin.RouterGroup) {
	middleware := middlewares.AuthMiddleware(oc.UserUseCase)
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
// @Failure 409,422,500
// @Failure 400,401 {object} models.DefaultErrorResponse
// @Router /user/orders [post]
func (oc OrdersController) orderCreate(c *gin.Context) {
	b, err := c.GetRawData()

	if err != nil {
		controllerLog(c).WithError(err).Errorln("order parse error")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	orderNr := string(b)
	token := c.GetHeader("Authorization")
	user, err := oc.UserUseCase.GetUser(token)

	if err != nil {
		controllerLog(c).WithError(err).Errorln("get user error")
		return
	}

	err = goluhn.Validate(orderNr)
	if err != nil {
		controllerLog(c).WithError(err).Errorln("order number is not valid")
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	orderRepo := oc.OrderUseCase.GetOrderRepository()
	hasOrder, err := orderRepo.HasOrder(orderNr)

	if err != nil {
		controllerLog(c).WithError(err).Errorln("order not found")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if hasOrder {
		order, err := orderRepo.GetOrder(orderNr)

		if err != nil {
			controllerLog(c).WithError(err).Errorln("order not found")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if order.UserID != user.UserID {
			c.AbortWithStatus(http.StatusConflict)
			return
		}

		c.AbortWithStatus(http.StatusOK)
		return
	}

	_, err = oc.OrderUseCase.CreateOrder(user.Login, orderNr)

	if err != nil {
		controllerLog(c).WithError(err).Errorln("can't create order")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	go func() {
		oc.AccrualService.ProcessOrder(orderNr)
	}()
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
// @Success 200 {json} []models.OrderListItem
// @Success 204
// @Failure 401 {object} models.DefaultErrorResponse
// @Failure 500
// @Router /user/orders [get]
func (oc OrdersController) orderGet(c *gin.Context) {
	token := c.GetHeader("Authorization")
	user, err := oc.UserUseCase.GetUser(token)

	if err != nil {
		controllerLog(c).WithError(err).Errorln("get user error")
		return
	}

	ordersList, err := oc.OrderUseCase.GetOrdersList(user.Login)
	if err != nil {
		controllerLog(c).WithError(err).Errorln("orders list error")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if len(ordersList) == 0 {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, ordersList)
}
