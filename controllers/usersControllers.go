package controllers

import (
	"motor/exceptions"
	"motor/models"
	"motor/payloads"
	"motor/repository"
	"motor/security"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var body struct {
		UserName  string `json:"username"`
		Email     string `json:"email"`
		Passsword string `json:"password"`
		Phone     string `json:"phone"`
		DeviceId  string `json:"device_id"`
	}

	c.ShouldBindJSON(&body)

	user := models.User{
		UserName: body.UserName,
		Email:    body.Email,
		Phone:    body.Phone,
		DeviceId: body.DeviceId,
	}

	findUserFromDB, _ := repository.GetUserByName(c, user.UserName)

	log.Info(findUserFromDB)

	if findUserFromDB.ID != 0 {
		exceptions.AppException(c, "User already exist")
		return
	}

	hash, err := security.HashPassword(body.Passsword)

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	user.Password = hash

	token, err := security.GenerateToken(user.UserName)

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	user.Token = token

	userResult, err := repository.CreateUser(c, user)

	if err != nil {
		exceptions.AppException(c, "Cant create user")
		return
	}

	member := models.Member{
		UserID: userResult.ID,
		Status: 0,
	}

	memberResult, _ := repository.CreateMember(c, member)

	if memberResult.Error != nil {
		exceptions.AppException(c, "Cant create member")
		return
	}

	payloads.HandleSuccess(c, userResult, "User created", 200)
}
