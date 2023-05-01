package repository

import (
	config "motor/configs"
	"motor/exceptions"
	"motor/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetMember(c *gin.Context, ID uint) (models.Member, error) {
	var member = models.Member{ID: ID}
	result := config.InitDB().First(&member)

	if result.Error != nil {
		exceptions.AppException(c, result.Error.Error())
		return member, result.Error
	}

	return member, nil

}

func CreateMember(c *gin.Context, member models.Member) (*gorm.DB, error) {
	result := config.InitDB().Create(&member)

	if result.Error != nil {
		exceptions.AppException(c, result.Error.Error())
		return nil, result.Error
	}

	return result, nil

}
