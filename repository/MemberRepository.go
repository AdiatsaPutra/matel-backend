package repository

import (
	config "motor/configs"
	"motor/exceptions"
	"motor/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllMember(c *gin.Context) ([]models.Member, error) {
	var members []models.Member
	result := config.InitDB().Find(&members)

	if result.Error != nil {
		exceptions.AppException(c, result.Error.Error())
		return members, result.Error
	}

	return members, nil

}

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

func UpdateMember(c *gin.Context, member models.Member) (models.Member, error) {
	var newMember = models.Member{}

	result := config.InitDB().Model(&member).Where("id = ?", member.ID).Update("status", member.Status)

	newMember = member

	if result.Error != nil {
		exceptions.AppException(c, result.Error.Error())
		return newMember, result.Error
	}

	return newMember, nil

}
