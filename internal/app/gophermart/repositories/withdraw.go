package repositories

import (
	"context"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
	"time"
)

type WithdrawRepository interface {
	CreateWithdraw(userID int, orderNr string, sum float32) (models.Withdraw, error)
	GetWithdraw(orderNr string) (models.Withdraw, error)
	GetWithdrawList(userID int) ([]models.Withdraw, error)
	HasWithdrawals(userID int) (bool, error)
}

func NewWithdrawRepository(conn *pgx.Conn) WithdrawRepository {
	return withdraw{DBConn: conn}
}

type withdraw struct {
	DBConn *pgx.Conn
}

func (w withdraw) CreateWithdraw(
	userID int,
	orderNr string,
	sum float32,
) (models.Withdraw, error) {
	ctx := context.Background()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.
		Update("withdrawals").
		Set("order_id", orderNr).
		Set("user_id", userID).
		Set("processed_at", time.Now()).
		Set("sum", sum).
		ToSql()

	if err != nil {
		log.WithFields(log.Fields{"func": "CreateWithdraw"}).Errorln(err)
		return models.Withdraw{}, err
	}

	_, err = w.DBConn.Exec(ctx, sql, args...)

	if err != nil {
		log.WithFields(log.Fields{"func": "CreateWithdraw"}).Errorln(err)
		return models.Withdraw{}, err
	}

	withdraw, err := w.GetWithdraw(orderNr)

	if err != nil {
		log.WithFields(log.Fields{"func": "CreateWithdraw"}).Errorln(err)
		return models.Withdraw{}, err
	}

	return withdraw, nil
}

func (w withdraw) GetWithdraw(orderNr string) (models.Withdraw, error) {
	ctx := context.Background()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.
		Select().
		Columns("id", "user_id", "order_id", "sum", "processed_at", "created_at").
		From("withdrawals").
		Where(sq.Eq{"order_id": orderNr}).
		ToSql()

	if err != nil {
		log.WithFields(log.Fields{"func": "GetWithdraw"}).Errorln(err)
		return models.Withdraw{}, err
	}

	rows, err := w.DBConn.Query(ctx, sql, args...)

	if err != nil {
		log.WithFields(log.Fields{"func": "GetWithdraw"}).Errorln(err)
		return models.Withdraw{}, err
	}

	defer rows.Close()

	var wd models.Withdraw
	rows.Next()

	err = rows.Scan(
		&wd.ID,
		&wd.UserID,
		&wd.Order,
		&wd.Sum,
		&wd.ProcessedAt,
		&wd.CreatedAt)

	if err != nil {
		log.WithFields(log.Fields{"func": "GetWithdraw"}).Errorln(err)
		return models.Withdraw{}, err
	}

	return wd, nil
}

func (w withdraw) GetWithdrawList(userID int) ([]models.Withdraw, error) {
	ctx := context.Background()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.
		Select().
		Columns("id", "user_id", "order_id", "sum", "processed_at", "created_at").
		From("withdrawals").
		Where(sq.Eq{"user_id": userID}).
		ToSql()

	if err != nil {
		log.WithFields(log.Fields{"func": "GetWithdrawList"}).Errorln(err)
		return []models.Withdraw{}, err
	}

	rows, err := w.DBConn.Query(ctx, sql, args...)

	if err != nil {
		log.WithFields(log.Fields{"func": "GetWithdrawList"}).Errorln(err)
		return []models.Withdraw{}, err
	}

	defer rows.Close()
	var wds []models.Withdraw
	for rows.Next() {
		var item models.Withdraw
		err = rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Order,
			&item.Sum,
			&item.ProcessedAt,
			&item.CreatedAt)

		if err != nil {
			log.WithFields(log.Fields{"func": "GetWithdrawList"}).Errorln(err)
			return []models.Withdraw{}, err
		}

		wds = append(wds, item)
	}

	return wds, nil
}

func (w withdraw) HasWithdrawals(userID int) (bool, error) {
	ctx := context.Background()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.
		Select("COALESCE(COUNT(*), 0) as count").
		From("withdrawals").
		Where(sq.Eq{"user_id": userID}).
		ToSql()

	if err != nil {
		log.WithFields(log.Fields{"func": "HasWithdrawals"}).Errorln(err)
		return false, err
	}

	rows, err := w.DBConn.Query(ctx, sql, args...)

	if err != nil {
		log.WithFields(log.Fields{"func": "HasWithdrawals"}).Errorln(err)
		return false, err
	}

	defer rows.Close()

	var count int
	rows.Next()
	err = rows.Scan(&count)

	if err != nil {
		log.WithFields(log.Fields{"func": "HasWithdrawals"}).Errorln(err)
		return false, err
	}

	return count > 0, nil
}
