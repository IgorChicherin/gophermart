package usecases

import (
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/models"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/repositories"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

type WithdrawUseCase interface {
	WithdrawalsList(userID int) ([]models.Withdraw, error)
}

func NewWithdrawUseCase(conn *pgx.Conn, withdrawRepo repositories.WithdrawRepository) WithdrawUseCase {
	return withdrawUseCase{DBConn: conn, WithdrawRepository: withdrawRepo}
}

type withdrawUseCase struct {
	DBConn             *pgx.Conn
	WithdrawRepository repositories.WithdrawRepository
}

func (c withdrawUseCase) WithdrawalsList(userID int) ([]models.Withdraw, error) {
	hasWd, err := c.WithdrawRepository.HasWithdrawals(userID)

	if err != nil {
		log.WithFields(log.Fields{"func": "WithdrawalsList"}).Errorln(err)
		return []models.Withdraw{}, err
	}

	if !hasWd {
		return []models.Withdraw{}, nil
	}

	wds, err := c.WithdrawRepository.GetWithdrawList(userID)

	if err != nil {
		log.WithFields(log.Fields{"func": "WithdrawalsList"}).Errorln(err)
		return []models.Withdraw{}, err
	}

	return wds, nil
}
