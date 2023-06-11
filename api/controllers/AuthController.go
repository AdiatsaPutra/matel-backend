package controllers

import (
	"matel/exceptions"
	"matel/models"
	"matel/payloads"
	"matel/repository"
	"matel/security"
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

	findUserFromDB, _ := repository.GetUserByEmail(c, body.Email)

	if findUserFromDB.ID != 0 {
		exceptions.AppException(c, "User sudah terdaftar")
		return
	}

	findUserDeviceID, _ := repository.GetUserByDeviceID(c, body.DeviceID)

	if findUserDeviceID.DeviceID == body.DeviceID {
		exceptions.AppException(c, "Perangkat anda sudah terdaftar")
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

	u, _ := repository.GetUserByEmail(c, userResult.Email)

	if u.UserName != "" {

		token, err := security.GenerateToken(u.ID)

		if err != nil {
			exceptions.AppException(c, err.Error())
			return
		}

		u.Token = token

		tokenCreated, err := repository.SetToken(c, u)

		if !tokenCreated {
			exceptions.AppException(c, err.Error())
			return
		}

		payloads.HandleSuccess(c, u, "Register Success", http.StatusOK)
	} else {
		exceptions.AppException(c, "Wrong Data")
		return
	}
}

func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		DeviceID string `json:"device_id"`
	}

	c.ShouldBindJSON(&body)

	findUserFromDB, _ := repository.GetUserByEmail(c, body.Email)

	if findUserFromDB.UserName != "" {

		if findUserFromDB.DeviceID == body.DeviceID || (findUserFromDB.DeviceID != body.DeviceID && findUserFromDB.Token == "") {
			hashPwd := findUserFromDB.Password
			pwd := body.Password

			hash := security.VerifyPassword(hashPwd, pwd)

			if hash == nil {
				if findUserFromDB.DeviceID == ""{
					err := repository.ResetDeviceID(c, findUserFromDB.ID, body.DeviceID)

					if err != nil {
						exceptions.AppException(c, err.Error())
						return
					}
				}

				if findUserFromDB.Token == "" {

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
				}
				payloads.HandleSuccess(c, findUserFromDB, "Login Success", http.StatusOK)
			} else {
				exceptions.AppException(c, "Wrong Data")
				return
			}
		} else {
			exceptions.AppException(c, "Data anda telah login")
			return
		}
	} else {
		exceptions.AppException(c, "User not registered")
		return
	}
}

func LoginWeb(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	c.ShouldBindJSON(&body)

	findUserFromDB, _ := repository.GetUserByEmail(c, body.Email)

	if findUserFromDB.UserName != "" {

		hashPwd := findUserFromDB.Password
		pwd := body.Password

		hash := security.VerifyPassword(hashPwd, pwd)

		if hash == nil {
			if findUserFromDB.Token == "" {

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

				c.JSON(http.StatusOK, gin.H{"token": findUserFromDB.Token})
			}
			c.JSON(http.StatusOK, gin.H{"token": findUserFromDB.Token})
		} else {
			exceptions.AppException(c, "Wrong Data")
			return
		}
	} else {
		exceptions.AppException(c, "User not registered")
		return
	}
}

func Logout(c *gin.Context) {
	UserID := c.MustGet("user_id").(uint)

	err := repository.Logout(c, UserID)

	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	payloads.HandleSuccess(c, "Success logout", "Success", 200)
}

func ResetDeviceID(c *gin.Context) {
	UserID := c.MustGet("user_id").(uint)

	err := repository.ResetDeviceID(c, UserID, "")

	if err != nil {
		exceptions.AppException(c, "Something went wrong")
		return
	}

	payloads.HandleSuccess(c, "Success reset", "Success", 200)
}
