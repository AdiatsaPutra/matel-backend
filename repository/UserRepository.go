package repository

import (
	config "motor/configs"
	"motor/exceptions"
	"motor/models"

	"github.com/gin-gonic/gin"
)

func GetUserByName(c *gin.Context, UserEmail string) (models.User, error) {
	var user = models.User{Email: UserEmail}
	result := config.InitDB().Where("email = ?", user.Email).First(&user)

	if result.Error != nil {
		return user, result.Error
	}

	return user, nil

}

func SetToken(c *gin.Context, user models.User) (bool, error) {
	result := config.InitDB().Model(&user).Where("id = ?", user.ID).Update("token", user.Token)

	if result.Error != nil {
		return true, result.Error
	}

	return true, nil

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

func UserProfile(c *gin.Context, user models.User) (models.User, error) {
	var newUser = models.User{}
	result := config.InitDB().First(&user)

	newUser = user

	if result.Error != nil {
		exceptions.AppException(c, result.Error.Error())
		return newUser, result.Error
	}

	return newUser, nil
}
