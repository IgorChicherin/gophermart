package controllers

import (
	"errors"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/models"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/repositories"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func controllerLog(c *gin.Context) *log.Entry {
	entryRaw, ok := c.Get("logger")
	if !ok {
		return log.NewEntry(log.StandardLogger())
	}

	entry, ok := entryRaw.(*log.Entry)
	if !ok {
		return log.NewEntry(log.StandardLogger())
	}

	return entry
}

func GetUser(c *gin.Context, userRepo repositories.UserRepository) (models.User, error) {
	token := c.GetHeader("Authorization")

	if token == "" {
		controllerLog(c).Errorln("unauthorized")
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized user"})
		return models.User{}, errors.New("unauthorized")
	}

	login, _, err := userRepo.DecodeToken(token)

	if err != nil {
		controllerLog(c).WithError(err).Errorln("can't decode token")
		c.AbortWithStatus(http.StatusInternalServerError)
		return models.User{}, errors.New("can't decode token")
	}

	user, err := userRepo.GetUser(login)

	if err != nil {
		controllerLog(c).WithError(err).Errorln("getting user error")
		c.AbortWithStatus(http.StatusInternalServerError)
		return models.User{}, errors.New("getting user error")
	}
	return user, nil
}
