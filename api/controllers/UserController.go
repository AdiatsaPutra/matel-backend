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

	var user models.User

	user.Status = newUser.Status
	user.StartSubscription = newUser.StartSubscription
	user.EndSubscription = newUser.EndSubscription
	user.CreatedAt = newUser.CreatedAt

	newUser.Status = uint(helper.GetUserStatus(user))

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
