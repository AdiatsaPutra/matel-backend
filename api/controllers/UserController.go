package controllers

import (
	config "matel/configs"
	"matel/exceptions"
	"matel/helper"
	"matel/models"
	"matel/payloads"
	"matel/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

	var user models.User

	user.StartSubscription = newUser.StartSubscription
	user.EndSubscription = newUser.EndSubscription
	user.CreatedAt = newUser.CreatedAt

	logrus.Info(newUser.SubscriptionMonth)
	logrus.Info(newUser.StartSubscription)
	logrus.Info(newUser.EndSubscription)

	newUser.Status = uint(helper.GetUserStatus(user))
	user.Status = newUser.Status

	if user.Status == 0 {
		var endDate = newUser.CreatedAt.Add(1 * 24 * time.Hour)
		user.EndSubscription = endDate.Format("2006-01-02")
		newUser.EndSubscription = user.EndSubscription
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

	var newUser []models.User

	user, err := repository.GetMember(c, search)

	for _, v := range user {
		v.Status = uint(helper.GetUserStatus(v))
		newUser = append(newUser, v)
	}

	for _, v := range newUser {
		logrus.Info(v.SubscriptionMonth)
	}

	if len(user) == 0 {
		payloads.HandleSuccess(c, nil, "User tidak ditemukan", http.StatusOK)
		return
	}

	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	payloads.HandleSuccess(c, newUser, "Success get data", http.StatusOK)
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

	logrus.Info("-------")
	logrus.Info(sub)

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

	err := config.InitDB().Unscoped().Delete(&user).Error
	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}
	payloads.HandleSuccess(c, "Berhasil mengubah status", "Berhasil", http.StatusOK)
}
