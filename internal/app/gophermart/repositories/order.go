package repositories

import (
	"context"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

type OrderRepository interface {
	CreateOrder(orderNr string, userID int) (models.Order, error)
	GetOrder(orderNr string) (models.Order, error)
	GetOrderList(userID int) ([]models.Order, error)
	HasOrder(orderNr string) (bool, error)
}

func NewOrderRepository(conn *pgx.Conn) OrderRepository {
	return orderRepo{DBConn: conn}
}

type orderRepo struct {
	DBConn *pgx.Conn
}

func (or orderRepo) CreateOrder(orderNr string, userID int) (models.Order, error) {
	ctx := context.Background()
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := psql.
		Insert("orders").
		Columns("order_id", "user_id").
		Values(orderNr, userID)
	sql, args, err := query.ToSql()

	if err != nil {
		log.Errorln(err)
		return models.Order{}, err
	}

	_, err = or.DBConn.Exec(ctx, sql, args...)
	if err != nil {
		log.Errorln(err)
		return models.Order{}, err
	}

	order, err := or.GetOrder(orderNr)

	if err != nil {
		log.Errorln(err)
		return models.Order{}, err
	}

	return order, nil
}

func (or orderRepo) GetOrder(orderNr string) (models.Order, error) {
	ctx := context.Background()
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Select().
		Columns("id", "order_id", "user_id", "status", "updated_at", "created_at").
		From("orders").
		Where(sq.Eq{"order_id": orderNr}).
		ToSql()

	if err != nil {
		log.Errorln(err)
		return models.Order{}, err
	}

	rows, err := or.DBConn.Query(ctx, sql, args...)
	if err != nil {
		log.Errorln(err)
		return models.Order{}, err
	}

	defer rows.Close()

	var order models.Order

	rows.Next()
	err = rows.Scan(&order.ID, &order.OrderID, &order.UserID, &order.Status, &order.UpdatedAt, &order.CreatedAt)

	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (or orderRepo) HasOrder(orderNr string) (bool, error) {
	ctx := context.Background()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.
		Select("COUNT(*)").
		From("orders").
		Where(sq.Eq{"order_id": orderNr}).
		ToSql()

	if err != nil {
		log.Errorln(err)
		return false, err
	}

	rows, err := or.DBConn.Query(ctx, sql, args...)

	if err != nil {
		log.Errorln(err)
		return false, err
	}

	defer rows.Close()

	var count int
	rows.Next()
	err = rows.Scan(&count)

	if err != nil {
		log.Errorln(err)
		return false, err
	}

	return count > 0, nil
}

func (or orderRepo) GetOrderList(userID int) ([]models.Order, error) {
	ctx := context.Background()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.
		Select().
		Columns("id", "order_id", "user_id", "status", "updated_at", "created_at").
		From("orders").
		Where(sq.Eq{"user_id": userID}).
		ToSql()

	if err != nil {
		log.Errorln(err)
		return []models.Order{}, err
	}

	rows, err := or.DBConn.Query(ctx, sql, args...)
	defer rows.Close()

	var ordersList []models.Order

	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.ID, &order.OrderID, &order.UserID, &order.Status, &order.UpdatedAt, &order.CreatedAt)

		if err != nil {
			log.Errorln(err)
			return []models.Order{}, err
		}

		ordersList = append(ordersList, order)
	}

	return ordersList, nil
}
