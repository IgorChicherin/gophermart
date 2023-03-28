package usecases

import (
	"context"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/models"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/repositories"
	"github.com/IgorChicherin/gophermart/internal/pkg/moneylib"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

type BalanceUseCase interface {
	GetBalance(login string) (models.Balance, error)
}

func NewBalanceUseCase(
	conn *pgx.Conn,
	userRepo repositories.UserRepository,
	moneyService moneylib.MoneyService,
) BalanceUseCase {
	return balance{DBConn: conn, UserRepo: userRepo, MoneyService: moneyService}
}

type balance struct {
	DBConn       *pgx.Conn
	UserRepo     repositories.UserRepository
	MoneyService moneylib.MoneyService
}

func (b balance) GetBalance(login string) (models.Balance, error) {
	user, err := b.UserRepo.GetUser(login)

	if err != nil {
		log.WithFields(log.Fields{"func": "GetBalance"}).Errorln(err)
		return models.Balance{}, err
	}

	ctx := context.Background()
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, err := psql.Select("COALESCE(SUM(accrual), 0) AS accrual, COALESCE(SUM(sum), 0) AS sum").
		From("orders").
		Where(sq.Eq{"user_id": user.UserID, "status": repositories.StatusProcessed}).
		LeftJoin("withdrawals USING (user_id)").
		ToSql()

	if err != nil {
		log.WithFields(log.Fields{"func": "GetBalance"}).Errorln(err)
		return models.Balance{}, err
	}

	rows, err := b.DBConn.Query(ctx, sql, args...)

	if err != nil {
		log.WithFields(log.Fields{"func": "GetBalance"}).Errorln(err)
		return models.Balance{}, err
	}

	defer rows.Close()

	var accrual, withdrawn int
	var balance models.Balance

	rows.Next()
	err = rows.Scan(&accrual, &withdrawn)

	if err != nil {
		log.WithFields(log.Fields{"func": "GetBalance"}).Errorln(err)
		return models.Balance{}, err
	}

	current := accrual - withdrawn

	balance.Current = b.MoneyService.IntToFloat32(current)
	balance.Withdrawn = b.MoneyService.IntToFloat32(withdrawn)
	return balance, nil
}
