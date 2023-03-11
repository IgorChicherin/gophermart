package router

import (
	docs "github.com/IgorChicherin/gophermart/api"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/controllers"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.RedirectTrailingSlash = false

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	docs.SwaggerInfo.BasePath = "/api/"

	auth := new(controllers.AuthController)
	orders := new(controllers.OrdersController)
	balance := new(controllers.BalanceController)
	withdraw := new(controllers.WithdrawController)

	api := router.Group("/api")
	{
		auth.Route(api)
		orders.Route(api)
		balance.Route(api)
		withdraw.Route(api)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}
