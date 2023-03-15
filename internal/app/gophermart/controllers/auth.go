package controllers

import (
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/models"
	"github.com/IgorChicherin/gophermart/internal/app/gophermart/repositories"
	"github.com/IgorChicherin/gophermart/internal/pkg/authlib"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type AuthController struct {
	UserRepository repositories.UserRepository
	AuthService    authlib.AuthService
}

func (ac AuthController) Route(api *gin.RouterGroup) {
	auth := api.Group("/user")
	{
		auth.POST("/login", ac.login)
		auth.POST("/register", ac.register)
	}
}

// @BasePath /api
// login godoc
// @Summary login
// @Schemes
// @Description user login
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.Login true "login"
// @Success 200
// @Failure 400,401,500
// @Router /user/login [post]
func (ac AuthController) login(c *gin.Context) {
	var userData models.Login

	if err := c.ShouldBind(&userData); err != nil {
		log.Errorln(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	hasLogin, err := ac.UserRepository.HasLogin(userData.Login)

	if err != nil {
		log.Errorln(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if !hasLogin {
		c.Status(http.StatusUnauthorized)
		return
	}

	user, err := ac.UserRepository.GetUser(userData.Login)

	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if !ac.AuthService.Equals(user.Password, userData.Password) {
		c.Status(http.StatusUnauthorized)
		return
	}

	token := ac.AuthService.EncodeToken(user.Login, user.Password)
	c.SetCookie("token", token, 3600, "/", "localhost", false, true)
	c.Status(http.StatusOK)
}

// @BasePath /api
// login godoc
// @Summary register
// @Schemes
// @Description user registration
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.User true "user account"
// @Success 200
// @Failure 404,500
// @Failure 400,409 {object} models.DefaultErrorResponse
// @Router /user/register [post]
func (ac AuthController) register(c *gin.Context) {
	var userData models.User

	if err := c.ShouldBind(&userData); err != nil {
		log.Errorln(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	hasLogin, err := ac.UserRepository.HasLogin(userData.Login)

	if err != nil {
		log.Errorln(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if hasLogin {
		c.AbortWithStatusJSON(http.StatusConflict, map[string]string{"error": "user with this login has been created"})
		return
	}

	createdUser, err := ac.UserRepository.CreateUser(userData.Login, userData.Password)

	if err != nil {
		log.Errorln(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	token := ac.AuthService.EncodeToken(createdUser.Login, createdUser.Password)
	c.SetCookie("token", token, 3600, "/", "localhost", false, true)

	c.Status(http.StatusOK)
}
