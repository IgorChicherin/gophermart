package usecases

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/IgorChicherin/gophermart/internal/app/gophermart/models"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/repositories"
	"github.com/IgorChicherin/gophermart/internal/pkg/authlib"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrUnauthorized          = errors.New("unauthorized user")
	ErrUserHasBeenRegistered = errors.New("user with this login has been created")
)

type UserUseCase interface {
	Validate(hash string) (bool, error)
	GetUser(token string) (models.User, error)
	Login(login, password string) (string, error)
	Register(login, password string) (string, error)
}

func NewUserUseCase(
	authService authlib.AuthService,
	userRepository repositories.UserRepository,
) UserUseCase {
	return userUseCase{AuthService: authService, UserRepository: userRepository}
}

type userUseCase struct {
	AuthService    authlib.AuthService
	UserRepository repositories.UserRepository
}

func (u userUseCase) Validate(hash string) (bool, error) {
	login, hash, err := u.AuthService.DecodeToken(hash)

	if err != nil {
		log.WithFields(log.Fields{"func": "Validate"}).Errorln(err)
		return false, err
	}

	user, err := u.UserRepository.GetUser(login)
	if err != nil {
		return false, err
	}
	return user.Password == hash, nil
}

func (u userUseCase) GetUser(token string) (models.User, error) {
	login, hash, err := u.AuthService.DecodeToken(token)

	if err != nil {
		return models.User{}, err
	}

	isCorrectPwd, err := u.Validate(hash)

	if err != nil {
		return models.User{}, err
	}

	if !isCorrectPwd {
		return models.User{}, errors.New("user password incorrect")
	}

	return u.UserRepository.GetUser(login)
}

func (u userUseCase) Login(login, password string) (string, error) {
	hasLogin, err := u.UserRepository.HasLogin(login)

	if !hasLogin {
		return "", ErrUserNotFound
	}

	if err != nil {
		return "", err
	}

	user, err := u.UserRepository.GetUser(login)

	if err != nil {
		return "", err
	}

	if !u.AuthService.Equals(user.Password, password) {
		return "", ErrUnauthorized
	}

	return u.AuthService.EncodeToken(login, password), nil
}

func (u userUseCase) Register(login, password string) (string, error) {
	hasLogin, err := u.UserRepository.HasLogin(login)

	if hasLogin {
		return "", ErrUserHasBeenRegistered
	}

	if err != nil {
		return "", err
	}

	createdUser, err := u.UserRepository.CreateUser(login, password)

	if err != nil {
		return "", err
	}

	return u.AuthService.EncodeToken(createdUser.Login, createdUser.Password), nil
}
