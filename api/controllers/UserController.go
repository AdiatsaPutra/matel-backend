package controllers

import (
	"matel/exceptions"
	"matel/payloads"
	"matel/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	UserID := c.MustGet("user_id").(uint)

	if UserID == 0 {
		exceptions.AppException(c, "Not authorized")
		return
	}

	newUser, err := repository.UserProfile(c, UserID)

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
		targetDate := newUser.CreatedAt.AddDate(0, 0, int((newUser.SubscriptionMonth*30)-1))

		currentDate := time.Now()

		var status = 0

		if currentDate.After(targetDate) {
			status = 2
			newUser.Status = uint(status)
		}

	}

	payloads.HandleSuccess(c, newUser, "Success get data", http.StatusOK)
}

func GetMember(c *gin.Context) {
	UserID := c.MustGet("user_id").(uint)
	search := c.Query("search")

	if UserID == 0 {
		exceptions.AppException(c, "Not authorized")
		return
	}

	user, err := repository.GetMember(c, search)

	if len(user) == 0 {
		payloads.HandleSuccess(c, nil, "User tidak ditemukan", http.StatusOK)
		return
	}

	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	payloads.HandleSuccess(c, user, "Success get data", http.StatusOK)
}

func SetUser(c *gin.Context) {
	type SetUserReq struct {
		UserID            uint `json:"user_id" validate:"required"`
		SubscriptionMonth uint `json:"subscription_month" validate:"required"`
	}
	var req SetUserReq
	c.BindJSON(&req)

	err := repository.SetUser(c, req.UserID, req.SubscriptionMonth)

	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	payloads.HandleSuccess(c, "Berhasil mengubah status", "Berhasil", http.StatusOK)
}
