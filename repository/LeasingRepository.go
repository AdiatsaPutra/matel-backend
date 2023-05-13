package repository

import (
	config "motor/configs"
	"motor/exceptions"
	"motor/models"

	"github.com/gin-gonic/gin"
)

func GetLeasing(c *gin.Context) ([]models.Leasing, error) {
	var leasing []models.Leasing
	result := config.InitDB().Find(&leasing)

	if result.Error != nil {
		exceptions.AppException(c, result.Error.Error())
		return leasing, result.Error
	}

	return leasing, nil

}
