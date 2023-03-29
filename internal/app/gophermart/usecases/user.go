package usecases

import (
	"errors"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/models"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/repositories"
	"github.com/IgorChicherin/gophermart/internal/pkg/authlib"
	log "github.com/sirupsen/logrus"
)

type UserUseCase interface {
	Validate(hash string) (bool, error)
	GetUser(token string) (models.User, error)
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
