package router

import (
	docs "github.com/IgorChicherin/gophermart/api"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/controllers"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/repositories"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/usecases"
	"github.com/IgorChicherin/gophermart/internal/pkg/accrual"
	"github.com/IgorChicherin/gophermart/internal/pkg/authlib"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(
	conn *pgx.Conn,
	authService authlib.AuthService,
	accrualService accrual.AccrualService,
) *gin.Engine {
	router := gin.New()
	router.RedirectTrailingSlash = false

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	docs.SwaggerInfo.BasePath = "/api/"

	userRepo := repositories.NewUserRepository(conn, authService)
	orderRepo := repositories.NewOrderRepository(conn)
	withdrawRepo := repositories.NewWithdrawRepository(conn)
	orderControllerUseCase := usecases.NewCreateOrderUseCase(conn, userRepo, orderRepo)
	balanceControllerUseCase := usecases.NewBalanceUseCase(conn, userRepo)
	withdrawUseCase := usecases.NewWithdrawUseCase(conn, withdrawRepo, balanceControllerUseCase)

	auth := controllers.AuthController{UserRepository: userRepo, AuthService: authService}
	orders := controllers.OrdersController{OrderUseCase: orderControllerUseCase, AccrualService: accrualService}
	balance := controllers.BalanceController{UserRepository: userRepo, BalanceUseCase: balanceControllerUseCase}
	withdraw := controllers.WithdrawController{UserRepository: userRepo, WithdrawUseCase: withdrawUseCase}

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
