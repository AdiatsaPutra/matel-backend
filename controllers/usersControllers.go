package controllers

import (
	"motor/configs"
	"motor/exception"
	"motor/models"
	"motor/payload"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var body struct {
		UserName string `json:"username"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		DeviceId string `json:"device_id"`
	}

	c.Bind(&body)

	user := models.User{
		UserName: body.UserName,
		Email:    body.Email,
		Phone:    body.Phone,
		DeviceId: body.DeviceId,
	}

	userResult := configs.DB.Create(&user)

	if userResult.Error != nil {
		exception.AppException(c, "Cant create user")
	}

	member := models.Member{
		UserID: user.ID,
		Status: 0,
	}

	memberResult := configs.DB.Create(&member)

	if memberResult.Error != nil {
		exception.AppException(c, "Cant create member")
	}

	payload.HandleSuccess(c, user, "Post created", 200)
}

func GetUser(c *gin.Context) {
	var post []models.User
	configs.DB.Find(&post)

	payload.HandleSuccess(c, post, "Post", 200)
}
