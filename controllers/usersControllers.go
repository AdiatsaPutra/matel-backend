package controllers

import (
	"motor/exceptions"
	"motor/models"
	"motor/payloads"
	"motor/repository"
	"motor/security"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var body struct {
		UserName    string `json:"username"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		Phone       string `json:"phone"`
		DeviceId    string `json:"device_id"`
		ProvinceID  int    `json:"province_id"`
		KabupatenID int    `json:"kabupaten_id"`
		KecamatanID int    `json:"kecamatan_id"`
	}

	c.ShouldBindJSON(&body)

	user := models.User{
		UserName: body.UserName,
		Email:    body.Email,
		Phone:    body.Phone,
	}

	findUserFromDB, _ := repository.GetUserByName(c, user.UserName)

	if findUserFromDB.ID != 0 {
		exceptions.AppException(c, "User already exist")
		return
	}

	hash, err := security.HashPassword(body.Password)

	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	user.Password = hash

	userResult, err := repository.CreateUser(c, user)

	if err != nil {
		exceptions.AppException(c, "Cant create user")
		return
	}

	member := models.Member{
		UserID:      userResult.ID,
		Status:      0,
		DeviceID:    body.DeviceID,
		ProvinceID:  body.ProvinceID,
		KabupatenID: body.KabupatenID,
		KecamatanID: body.KecamatanID,
	}

	memberResult, _ := repository.CreateMember(c, member)

	if memberResult.Error != nil {
		exceptions.AppException(c, "Cant create member")
		return
	}

	payloads.HandleSuccess(c, userResult, "User created", http.StatusOK)
}

func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	c.ShouldBindJSON(&body)

	findUserFromDB, _ := repository.GetUserByName(c, body.Email)

	if findUserFromDB.UserName != "" {

		hashPwd := findUserFromDB.Password
		pwd := body.Password

		hash := security.VerifyPassword(hashPwd, pwd)

		if hash == nil {
			token, err := security.GenerateToken(findUserFromDB.ID)

			if err != nil {
				exceptions.AppException(c, err.Error())
				return
			}

			findUserFromDB.Token = token

			tokenCreated, err := repository.SetToken(c, findUserFromDB)

			if !tokenCreated {
				exceptions.AppException(c, err.Error())
				return
			}

			payloads.HandleSuccess(c, findUserFromDB, "Login Success", http.StatusOK)
		} else {
			exceptions.AppException(c, "Wrong Data")
			return
		}
	} else {
		exceptions.AppException(c, "User not registered")
		return
	}
}
