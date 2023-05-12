package controllers

import (
	"motor/exceptions"
	"motor/models"
	"motor/payloads"
	"motor/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetProfile(c *gin.Context) {
	UserID := c.MustGet("user_id").(uint)

	logrus.Info(UserID)

	user := models.User{
		ID: UserID,
	}

	newUser, err := repository.UserProfile(c, user)

	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	payloads.HandleSuccess(c, newUser, "Success get data", http.StatusOK)
}