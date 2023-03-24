package controllers

import (
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

func GetUser(c *gin.Context, userRepo repositories.UserRepository) (error, models.User) {
	token := c.GetHeader("Authorization")

	if token == "" {
		controllerLog(c).Errorln("unauthorized")
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized user"})
		return nil, models.User{}
	}

	login, _, err := userRepo.DecodeToken(token)

	if err != nil {
		controllerLog(c).WithError(err).Errorln("can't decode token")
		c.AbortWithStatus(http.StatusInternalServerError)
		return nil, models.User{}
	}

	user, err := userRepo.GetUser(login)

	if err != nil {
		controllerLog(c).WithError(err).Errorln("getting user error")
		c.AbortWithStatus(http.StatusInternalServerError)
		return nil, models.User{}
	}
	return err, user
}
