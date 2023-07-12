package controllers

import (
	config "matel/configs"
	"matel/exceptions"
	"matel/models"
	"matel/payloads"
	"matel/repository"
	"net/http"
	"strconv"
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
			// status = 2
			newUser.Status = uint(status)
		}

		newUser.StartSubscription = newUser.CreatedAt.Format("2006-01-02 15:04:05")
		newUser.EndSubscription = targetDate.Format("2006-01-02 15:04:05")

	} else if newUser.Status == 1 {
		targetDate := newUser.CreatedAt.AddDate(0, 0, int((newUser.SubscriptionMonth*30)-1))

		currentDate := time.Now()

		var status = 0

		if currentDate.After(targetDate) {
			// status = 2
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
		UserID            uint   `json:"user_id" validate:"required"`
		SubscriptionMonth string `json:"subscription_month" validate:"required"`
	}
	var req SetUserReq
	c.BindJSON(&req)

	sub, e := strconv.Atoi(req.SubscriptionMonth)

	if e != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	err := repository.SetUser(c, req.UserID, uint(sub))

	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	payloads.HandleSuccess(c, "Berhasil mengubah status", "Berhasil", http.StatusOK)
}

func DeleteMember(c *gin.Context) {
	id := c.Param("id")

	user := models.User{}
	if err := config.InitDB().Where("id = ?", id).First(&user).Error; err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	err := config.InitDB().Delete(&user).Error
	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}
	payloads.HandleSuccess(c, "Berhasil mengubah status", "Berhasil", http.StatusOK)
}
