package controllers

import (
	"motor/exceptions"
	"motor/models"
	"motor/payloads"
	"motor/repository"
	"motor/security"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Register(c *gin.Context) {
	var body struct {
		UserName    string `json:"username" validate:"required"`
		Email       string `json:"email" validate:"required,email"`
		Password    string `json:"password" validate:"required"`
		Phone       string `json:"phone" validate:"required"`
		DeviceID    string `json:"device_id" validate:"required"`
		ProvinceID  uint   `json:"province_id" validate:"required"`
		KabupatenID uint   `json:"kabupaten_id" validate:"required"`
		KecamatanID uint   `json:"kecamatan_id" validate:"required"`
	}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	validate := validator.New()
	err = validate.Struct(body)
	if err != nil {
		exceptions.AppException(c, "Lengkapi data anda")
		return
	}

	findUserFromDB, err := repository.GetUserByName(c, body.UserName)
	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	emptyUser := models.User{}
	if findUserFromDB != emptyUser {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	hash, err := security.HashPassword(body.Password)
	if err != nil {
		exceptions.AppException(c, err.Error())
		return
	}

	user := models.User{
		UserName:    body.UserName,
		Email:       body.Email,
		Phone:       body.Phone,
		Status:      0,
		DeviceID:    body.DeviceID,
		ProvinceID:  body.ProvinceID,
		KabupatenID: body.KabupatenID,
		KecamatanID: body.KecamatanID,
		Password:    hash,
	}

	userResult, err := repository.CreateUser(c, user)
	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	payloads.HandleSuccess(c, userResult, "Success Register", http.StatusOK)
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
