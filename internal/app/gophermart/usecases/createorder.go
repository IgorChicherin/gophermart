package usecases

import (
	"errors"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/models"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/repositories"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

type CreateOrderUseCase interface {
	CreateOrder(login, orderNr string) (models.Order, error)
	GetUserRepository() repositories.UserRepository
	GetOrderRepository() repositories.OrderRepository
}

func NewCreateOrderUseCase(
	conn *pgx.Conn,
	userRepo repositories.UserRepository,
	orderRepo repositories.OrderRepository) CreateOrderUseCase {
	return createOrderUseCase{DBConn: conn, UserRepo: userRepo, OrderRepo: orderRepo}
}

type createOrderUseCase struct {
	DBConn    *pgx.Conn
	UserRepo  repositories.UserRepository
	OrderRepo repositories.OrderRepository
}

func (c createOrderUseCase) CreateOrder(login, orderNr string) (models.Order, error) {
	found, err := c.UserRepo.HasLogin(login)

	if err != nil {
		log.Errorln(err)
		return models.Order{}, err
	}

	if !found {
		return models.Order{}, errors.New("user not found")
	}

	user, err := c.UserRepo.GetUser(login)

	if err != nil {
		log.Errorln(err)
	}

	order, err := c.OrderRepo.CreateOrder(orderNr, user.UserId)

	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (c createOrderUseCase) GetUserRepository() repositories.UserRepository {
	return c.UserRepo
}

func (c createOrderUseCase) GetOrderRepository() repositories.OrderRepository {
	return c.OrderRepo
}
