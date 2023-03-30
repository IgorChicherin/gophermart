package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/models"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/repositories"
	"github.com/IgorChicherin/gophermart/internal/pkg/moneylib"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const checkOrderURL = "/api/orders/%s"

var (
	ErrNotFoundOrder = errors.New("order not found")
)

type OrderAccrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual"`
}

type AccrualService interface {
	Run()
	ProcessOrders() error
}

type accrual struct {
	Host         string
	DBConn       *pgx.Conn
	Ctx          context.Context
	MoneyService moneylib.MoneyService
	quitCh       chan struct{}
}

func NewAccrualService(
	ctx context.Context,
	conn *pgx.Conn,
	accrualHost string,
	moneyService moneylib.MoneyService,
) AccrualService {
	return accrual{
		Ctx:          ctx,
		DBConn:       conn,
		Host:         accrualHost,
		MoneyService: moneyService,
	}
}

func (a accrual) Run() {
	for {
		select {
		case <-a.Ctx.Done():
			log.Infof("accrual service shutting down")
			return
		default:
			err := a.ProcessOrders()
			if err != nil {
				log.WithFields(log.Fields{"func": "processOrders"}).Errorln(err)
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func (a accrual) ProcessOrders() error {
	orders, err := a.getOrders()

	if err != nil {
		log.WithFields(log.Fields{"func": "processOrders"}).Errorln(err)
		return err
	}

	for _, order := range orders {
		err := a.processOrder(order.OrderID)
		if err != nil {
			log.WithFields(log.Fields{"func": "processOrders"}).Errorln(err)
		}
	}
	return err
}

func (a accrual) getOrders() ([]models.Order, error) {
	ctx := context.Background()
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, err := psql.
		Select().
		Columns("id", "order_id", "user_id", "status", "updated_at", "created_at", "accrual").
		From("orders").
		Where(sq.Or{sq.Eq{"status": repositories.StatusNew}, sq.Eq{"status": repositories.StatusProcessing}}).
		ToSql()

	if err != nil {
		log.WithFields(log.Fields{"func": "getOrders"}).Errorln(err)
		return []models.Order{}, err
	}

	rows, err := a.DBConn.Query(ctx, sql, args...)

	if err != nil {
		log.WithFields(log.Fields{"func": "getOrders"}).Errorln(err)
		return []models.Order{}, err
	}

	defer rows.Close()

	var ordersList []models.Order

	for rows.Next() {
		var order models.Order

		err = rows.Scan(
			&order.ID,
			&order.OrderID,
			&order.UserID,
			&order.Status,
			&order.UpdatedAt,
			&order.CreatedAt,
			&order.Accrual)

		if err != nil {
			log.WithFields(log.Fields{"func": "getOrders"}).Errorln(err)
			return []models.Order{}, err
		}

		ordersList = append(ordersList, order)
	}
	return ordersList, nil
}

func (a accrual) getAccrual(orderNr string) (OrderAccrual, error) {
	URL := fmt.Sprintf(checkOrderURL, orderNr)

	req := resty.
		New().
		SetBaseURL(a.Host).
		R().
		SetHeader("Content-Type", "application/json")

	resp, err := req.Get(URL)

	if err != nil {
		log.WithFields(log.Fields{"func": "GetAccrual"}).Errorln(err)
		return OrderAccrual{}, err
	}

	if resp.StatusCode() == http.StatusNoContent {
		return OrderAccrual{}, ErrNotFoundOrder
	}

	data := resp.Body()
	var order OrderAccrual

	err = json.Unmarshal(data, &order)

	if err != nil {
		log.WithFields(log.Fields{"func": "GetAccrual"}).Errorln(err)
		return OrderAccrual{}, err
	}

	return order, nil
}

func (a accrual) processOrder(orderNr string) error {
	orderAccrual, err := a.getAccrual(orderNr)

	if errors.Is(err, ErrNotFoundOrder) {
		log.WithFields(log.Fields{"func": "processOrder"}).Errorln(err)
		return err
	}

	if err != nil {
		log.WithFields(log.Fields{"func": "processOrder"}).Errorln(err)
		return err
	}

	orderStatus := repositories.OrderStatus(orderAccrual.Status)

	if orderStatus == repositories.StatusProcessed || orderStatus == repositories.StatusInvalid {
		err = a.updateOrder(orderAccrual)
		if err != nil {
			log.WithFields(log.Fields{"func": "processOrder"}).Errorln(err)
			return err
		}
	}
	return nil
}

func (a accrual) updateOrder(order OrderAccrual) error {
	ctx := context.Background()
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Update("orders").
		Where(sq.Eq{"order_id": order.Order}).
		Set("status", order.Status).
		Set("accrual", a.MoneyService.FloatToInt(order.Accrual)).
		Set("updated_at", time.Now()).
		ToSql()

	if err != nil {
		log.WithFields(log.Fields{"func": "updateOrder"}).Errorln(err)
		return err
	}

	_, err = a.DBConn.Exec(ctx, sql, args...)

	if err != nil {
		log.WithFields(log.Fields{"func": "updateOrder"}).Errorln(err)
		return err
	}

	return nil
}
