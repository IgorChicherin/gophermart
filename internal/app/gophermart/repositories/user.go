package repositories

import (
	"context"
	"github.com/IgorChicherin/gophermart/internal/pkg/authlib"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

type UserRepository interface {
	GetUser(login string) (user, error)
	Validate(hash string) (bool, error)
}

type user struct {
	Login     string `json:"login"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

type userRepo struct {
	DBConn      *pgx.Conn
	AuthService authlib.AuthService
}

func NewUserRepository(conn *pgx.Conn, service authlib.AuthService) UserRepository {
	return userRepo{DBConn: conn, AuthService: service}
}

func (ur userRepo) GetUser(login string) (user, error) {
	ctx := context.Background()
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Select().
		Columns("login", "password", "CreatedAt").
		From("user").
		Where(sq.Eq{"login": login}).
		ToSql()

	if err != nil {
		log.Errorln(err)
		return user{}, err
	}

	rows, err := ur.DBConn.Query(ctx, sql, args...)
	if err != nil {
		log.Errorln(err)
		return user{}, err
	}

	defer rows.Close()

	var u user
	err = rows.Scan(&u.Login, &u.Password, &u.CreatedAt)

	if err != nil {
		return user{}, err
	}

	return u, nil
}

func (ur userRepo) Validate(hash string) (bool, error) {
	login, hash, err := ur.AuthService.DecodeToken(hash)

	user, err := ur.GetUser(login)
	if err != nil {
		return false, err
	}
	return user.Password == hash, nil
}
