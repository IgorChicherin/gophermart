package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/repositories"
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
	Accrual float64 `json:"accrual"`
}

type AccrualService interface {
	GetAccrual(orderNr string) (OrderAccrual, error)
	ProcessOrder(orderNr string)
}

type accrual struct {
	Host   string
	DBConn *pgx.Conn
	Ctx    context.Context
}

func NewAccrualService(ctx context.Context, conn *pgx.Conn, accrualHost string) AccrualService {
	return accrual{
		Ctx:    ctx,
		DBConn: conn,
		Host:   accrualHost,
	}
}

func (a accrual) GetAccrual(orderNr string) (OrderAccrual, error) {
	host := "http://" + a.Host
	URL := fmt.Sprintf(checkOrderURL, orderNr)

	req := resty.
		New().
		SetBaseURL(host).
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

func (a accrual) ProcessOrder(orderNr string) {
	for {
		select {
		case <-a.Ctx.Done():
			log.WithFields(log.Fields{"func": "ProcessOrder"}).
				Errorf("processing order #%s has been cancelled", orderNr)
			return
		default:
			orderAccrual, err := a.GetAccrual(orderNr)

			if errors.Is(err, ErrNotFoundOrder) {
				log.WithFields(log.Fields{"func": "ProcessOrder"}).Errorln(err)
				time.Sleep(1 * time.Second)
				continue
			}

			if err != nil {
				log.WithFields(log.Fields{"func": "ProcessOrder"}).Errorln(err)
				return
			}

			orderStatus := repositories.OrderStatus(orderAccrual.Status)

			if orderStatus == repositories.StatusProcessed || orderStatus == repositories.StatusInvalid {
				err = a.updateOrder(orderAccrual)
				if err != nil {
					log.WithFields(log.Fields{"func": "ProcessOrder"}).Errorln(err)
				}
				return
			}

			time.Sleep(1 * time.Second)
			continue
		}
	}
}

func (a accrual) updateOrder(order OrderAccrual) error {
	ctx := context.Background()
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Update("orders").
		Where(sq.Eq{"order_id": order.Order}).
		Set("status", order.Status).
		Set("accrual", order.Accrual).
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
