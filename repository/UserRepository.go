package repository

import (
	config "motor/configs"
	"motor/exceptions"
	"motor/models"

	"github.com/gin-gonic/gin"
)

func GetUserByName(c *gin.Context, UserName string) (models.User, error) {
	var user = models.User{UserName: UserName}
	result := config.InitDB().Where("user_name = ?", user.UserName).First(&user)

	if result.Error != nil {
		return user, result.Error
	}

	return user, nil

}

func CreateUser(c *gin.Context, user models.User) (models.User, error) {
	var newUser = models.User{}
	result := config.InitDB().Create(&user)

	newUser = user

	if result.Error != nil {
		exceptions.AppException(c, result.Error.Error())
		return newUser, result.Error
	}

	return newUser, nil

}
