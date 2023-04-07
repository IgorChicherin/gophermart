package usecases

import (
	"errors"

	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"

	"github.com/IgorChicherin/gophermart/internal/app/gophermart/models"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/repositories"
	"github.com/IgorChicherin/gophermart/internal/pkg/moneylib"
)

var ErrInsufficientFunds = errors.New("insufficient fund")

type WithdrawUseCase interface {
	WithdrawalsList(userID int) ([]models.WithdrawListItem, error)
	CreateWithdrawOrder(user models.User, orderNr string, sum float32) (models.Withdraw, error)
}

func NewWithdrawUseCase(
	conn *pgx.Conn,
	withdrawRepo repositories.WithdrawRepository,
	balanceUseCase BalanceUseCase,
	moneyService moneylib.MoneyService,
) WithdrawUseCase {
	return withdrawUseCase{
		DBConn:             conn,
		WithdrawRepository: withdrawRepo,
		BalanceUseCase:     balanceUseCase,
		MoneyService:       moneyService,
	}
}

type withdrawUseCase struct {
	DBConn             *pgx.Conn
	WithdrawRepository repositories.WithdrawRepository
	BalanceUseCase     BalanceUseCase
	MoneyService       moneylib.MoneyService
}

func (w withdrawUseCase) WithdrawalsList(userID int) ([]models.WithdrawListItem, error) {
	hasWd, err := w.WithdrawRepository.HasWithdrawals(userID)

	if err != nil {
		log.WithFields(log.Fields{"func": "WithdrawalsList"}).Errorln(err)
		return []models.WithdrawListItem{}, err
	}

	if !hasWd {
		return []models.WithdrawListItem{}, nil
	}

	wds, err := w.WithdrawRepository.GetWithdrawList(userID)

	if err != nil {
		log.WithFields(log.Fields{"func": "WithdrawalsList"}).Errorln(err)
		return []models.WithdrawListItem{}, err
	}

	var items []models.WithdrawListItem

	for _, item := range wds {
		items = append(items, models.WithdrawListItem{
			ID:          item.ID,
			Order:       item.Order,
			UserID:      item.UserID,
			Sum:         w.MoneyService.IntToFloat32(item.Sum),
			ProcessedAt: item.ProcessedAt,
			CreatedAt:   item.CreatedAt,
		})
	}

	return items, nil
}

func (w withdrawUseCase) CreateWithdrawOrder(user models.User, orderNr string, sum float32) (models.Withdraw, error) {

	balance, err := w.BalanceUseCase.GetBalance(user.Login)

	if err != nil {
		log.WithFields(log.Fields{"func": "CreateWithdrawOrder"}).Errorln(err)
		return models.Withdraw{}, err
	}

	if balance.Current < sum {
		log.WithFields(log.Fields{"func": "CreateWithdrawOrder"}).Errorln(ErrInsufficientFunds)
		return models.Withdraw{}, ErrInsufficientFunds
	}

	wdSum := w.MoneyService.FloatToInt(sum)
	return w.WithdrawRepository.CreateWithdraw(user.UserID, orderNr, wdSum)
}
