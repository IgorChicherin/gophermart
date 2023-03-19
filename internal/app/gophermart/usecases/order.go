package usecases

import (
	"errors"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/models"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/repositories"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

type OrderUseCase interface {
	GetUserRepository() repositories.UserRepository
	GetOrderRepository() repositories.OrderRepository
	CreateOrder(login, orderNr string) (models.Order, error)
	GetOrdersList(login string) ([]models.OrderListItem, error)
}

func NewCreateOrderUseCase(
	conn *pgx.Conn,
	userRepo repositories.UserRepository,
	orderRepo repositories.OrderRepository) OrderUseCase {
	return orderUseCase{DBConn: conn, UserRepo: userRepo, OrderRepo: orderRepo}
}

type orderUseCase struct {
	DBConn    *pgx.Conn
	UserRepo  repositories.UserRepository
	OrderRepo repositories.OrderRepository
}

func (c orderUseCase) CreateOrder(login, orderNr string) (models.Order, error) {
	found, err := c.UserRepo.HasLogin(login)

	if err != nil {
		log.WithFields(log.Fields{"func": "CreateOrder"}).Errorln(err)
		return models.Order{}, err
	}

	if !found {
		return models.Order{}, errors.New("user not found")
	}

	user, err := c.UserRepo.GetUser(login)

	if err != nil {
		log.WithFields(log.Fields{"func": "CreateOrder"}).Errorln(err)
	}

	order, err := c.OrderRepo.CreateOrder(orderNr, user.UserID)

	if err != nil {
		log.WithFields(log.Fields{"func": "CreateOrder"}).Errorln(err)
		return models.Order{}, err
	}

	return order, nil
}

func (c orderUseCase) GetOrdersList(login string) ([]models.OrderListItem, error) {
	var ordersListResponse []models.OrderListItem
	hasLogin, err := c.UserRepo.HasLogin(login)

	if err != nil {
		log.WithFields(log.Fields{"func": "GetOrdersList"}).Errorln(err)
		return ordersListResponse, err
	}

	if !hasLogin {
		return ordersListResponse, repositories.ErrUserNotFound
	}

	user, err := c.UserRepo.GetUser(login)

	if err != nil {
		log.WithFields(log.Fields{"func": "GetOrdersList"}).Errorln(err)
		return ordersListResponse, err
	}

	orders, err := c.OrderRepo.GetOrderList(user.UserID)

	if err != nil {
		log.WithFields(log.Fields{"func": "GetOrdersList"}).Errorln(err)
		return ordersListResponse, err
	}

	for _, order := range orders {
		ordersListResponse = append(ordersListResponse, models.OrderListItem{
			Number:     order.OrderID,
			Status:     order.Status,
			Accrual:    order.Accrual,
			UploadedAt: order.CreatedAt,
		})
	}
	return ordersListResponse, nil
}

func (c orderUseCase) GetUserRepository() repositories.UserRepository {
	return c.UserRepo
}

func (c orderUseCase) GetOrderRepository() repositories.OrderRepository {
	return c.OrderRepo
}
