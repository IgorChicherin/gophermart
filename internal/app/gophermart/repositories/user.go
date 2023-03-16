package repositories

import (
	"context"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/models"
	"github.com/IgorChicherin/gophermart/internal/pkg/authlib"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

type UserRepository interface {
	HasLogin(login string) (bool, error)
	GetUser(login string) (models.User, error)
	CreateUser(login, password string) (models.User, error)
	Validate(hash string) (bool, error)
	DecodeToken(token string) (string, string, error)
}

type userRepo struct {
	DBConn      *pgx.Conn
	AuthService authlib.AuthService
}

func NewUserRepository(conn *pgx.Conn, service authlib.AuthService) UserRepository {
	return userRepo{DBConn: conn, AuthService: service}
}

func (ur userRepo) GetUser(login string) (models.User, error) {
	ctx := context.Background()
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Select().
		Columns("id", "login", "password", "created_at").
		From("users").
		Where(sq.Eq{"login": login}).
		ToSql()

	if err != nil {
		log.Errorln(err)
		return models.User{}, err
	}

	rows, err := ur.DBConn.Query(ctx, sql, args...)
	if err != nil {
		log.Errorln(err)
		return models.User{}, err
	}

	defer rows.Close()

	var u models.User

	rows.Next()
	err = rows.Scan(&u.UserId, &u.Login, &u.Password, &u.CreatedAt)

	if err != nil {
		return models.User{}, err
	}

	return u, nil
}

func (ur userRepo) CreateUser(login, password string) (models.User, error) {
	ctx := context.Background()

	pwdHash := ur.AuthService.GetHash(password)

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query := psql.
		Insert("users").
		Columns("login", "password").
		Values(login, pwdHash)
	sql, args, err := query.ToSql()

	if err != nil {
		log.Errorln(err)
		return models.User{}, err
	}

	_, err = ur.DBConn.Exec(ctx, sql, args...)
	if err != nil {
		log.Errorln(err)
		return models.User{}, err
	}

	return models.User{Login: login, Password: pwdHash}, nil
}

func (ur userRepo) HasLogin(login string) (bool, error) {
	ctx := context.Background()

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.
		Select("COUNT(*)").
		From("users").
		Where(sq.Eq{"login": login}).
		ToSql()

	if err != nil {
		log.Errorln(err)
		return false, err
	}

	rows, err := ur.DBConn.Query(ctx, sql, args...)

	if err != nil {
		log.Errorln(err)
		return false, err
	}

	defer rows.Close()

	var count int
	rows.Next()
	err = rows.Scan(&count)

	return count > 0, nil
}

func (ur userRepo) Validate(hash string) (bool, error) {
	login, hash, err := ur.AuthService.DecodeToken(hash)

	user, err := ur.GetUser(login)
	if err != nil {
		return false, err
	}
	return user.Password == hash, nil
}

func (ur userRepo) DecodeToken(token string) (string, string, error) {
	return ur.AuthService.DecodeToken(token)
}
