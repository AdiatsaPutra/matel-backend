package repository

import (
	"database/sql"
	config "motor/configs"
	"motor/exceptions"

	"github.com/gin-gonic/gin"
)

func GetLeasingTotal(c *gin.Context) (uint, error) {
	var count sql.NullInt64
	result := config.InitDB().Raw("SELECT COUNT(*) FROM m_leasing").Scan(&count)

	if result.Error != nil {
		exceptions.AppException(c, result.Error.Error())
		return 0, result.Error
	}

	return uint(count.Int64), nil

}
