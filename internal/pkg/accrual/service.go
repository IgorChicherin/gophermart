package accrual

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"net/http"
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
}

type accrual struct {
	Host string
}

func NewAccrualService(accrualHost string) AccrualService {
	return accrual{Host: accrualHost}
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
