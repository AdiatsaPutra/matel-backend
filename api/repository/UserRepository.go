package repository

import (
	config "matel/configs"
	"matel/exceptions"
	"matel/models"

	"github.com/gin-gonic/gin"
)

func GetUserTotal(c *gin.Context) (uint, error) {
	var user []models.User
	result := config.InitDB().Find(&user)

	if result.Error != nil {
		return 0, result.Error
	}

	return uint(result.RowsAffected), nil

}

func GetUserByEmail(c *gin.Context, UserEmail string) (models.User, error) {
	var user = models.User{Email: UserEmail}
	result := config.InitDB().Where("email = ?", user.Email).First(&user)

	if result.Error != nil {
		return user, result.Error
	}

	return user, nil

}

func GetUserByDeviceID(c *gin.Context, DeviceID string) (models.User, error) {
	var user = models.User{}
	result := config.InitDB().Where("device_id = ?", DeviceID).First(&user)

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

func GetMember(c *gin.Context) ([]models.User, error) {
	var user = []models.User{}
	result := config.InitDB().Where("is_admin = 0").Find(&user)

	if result.Error != nil {
		exceptions.AppException(c, result.Error.Error())
		return user, result.Error
	}

	return user, nil
}

func Logout(c *gin.Context, UserID uint) error {
	var user models.User

	err := config.InitDB().Model(&user).Where("id = ?", UserID).Update("token", "").Error

	if err != nil {
		return err
	}
	return nil

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
