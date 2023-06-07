package controllers

import (
	"matel/exceptions"
	"matel/models"
	"matel/payloads"
	"matel/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	UserID := c.MustGet("user_id").(uint)

	user := models.User{
		ID: UserID,
	}

	newUser, err := repository.UserProfile(c, user)

	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	if newUser.Status == 0 {
		targetDate := newUser.CreatedAt.AddDate(0, 0, 29)

		currentDate := time.Now()

		var status = 0

		if currentDate.After(targetDate) {
			status = 2
			newUser.Status = uint(status)
		}

	} else if newUser.Status == 1 {
		targetDate := newUser.CreatedAt.AddDate(0, 0, int(newUser.SubscriptionMonth*30))

		currentDate := time.Now()

		var status = 0

		if currentDate.After(targetDate) {
			status = 2
			newUser.Status = uint(status)
		}

	}

	payloads.HandleSuccess(c, newUser, "Success get data", http.StatusOK)
}
