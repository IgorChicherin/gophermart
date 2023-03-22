package usecases

import (
	"errors"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/models"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/repositories"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

var ErrInsufficientFunds = errors.New("insufficient fund")

type WithdrawUseCase interface {
	WithdrawalsList(userID int) ([]models.Withdraw, error)
	CreateWithdrawOrder(user models.User, orderNr string, sum float32) (models.Withdraw, error)
}

func NewWithdrawUseCase(
	conn *pgx.Conn,
	withdrawRepo repositories.WithdrawRepository,
	balanceUseCase BalanceUseCase,
) WithdrawUseCase {
	return withdrawUseCase{DBConn: conn, WithdrawRepository: withdrawRepo, BalanceUseCase: balanceUseCase}
}

type withdrawUseCase struct {
	DBConn             *pgx.Conn
	WithdrawRepository repositories.WithdrawRepository
	BalanceUseCase     BalanceUseCase
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

func (c withdrawUseCase) CreateWithdrawOrder(user models.User, orderNr string, sum float32) (models.Withdraw, error) {

	balance, err := c.BalanceUseCase.GetBalance(user.Login)

	if err != nil {
		log.WithFields(log.Fields{"func": "CreateWithdrawOrder"}).Errorln(err)
		return models.Withdraw{}, err
	}

	if balance.Current < sum {
		log.WithFields(log.Fields{"func": "CreateWithdrawOrder"}).Errorln(ErrInsufficientFunds)
		return models.Withdraw{}, ErrInsufficientFunds
	}

	return c.WithdrawRepository.CreateWithdraw(user.UserID, orderNr, sum)
}
