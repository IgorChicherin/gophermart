package router

import (
	docs "github.com/IgorChicherin/gophermart/api"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/controllers"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/repositories"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/usecases"
	"github.com/IgorChicherin/gophermart/internal/pkg/authlib"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(conn *pgx.Conn, authService authlib.AuthService) *gin.Engine {
	router := gin.New()
	router.RedirectTrailingSlash = false

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	docs.SwaggerInfo.BasePath = "/api/"

	userRepo := repositories.NewUserRepository(conn, authService)
	orderRepo := repositories.NewOrderRepository(conn)
	orderControllerUseCase := usecases.NewCreateOrderUseCase(conn, userRepo, orderRepo)

	auth := controllers.AuthController{UserRepository: userRepo, AuthService: authService}
	orders := controllers.OrdersController{CreateOrderUseCase: orderControllerUseCase}
	balance := controllers.BalanceController{UserRepository: userRepo}
	withdraw := controllers.WithdrawController{UserRepository: userRepo}

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
