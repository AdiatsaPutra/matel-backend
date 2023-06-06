package exceptions

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppExceptionError struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func AppException(c *gin.Context, message string) {
	res := AppExceptionError{
		Success: false,
		Message: message,
		Data:    nil,
	}

	c.JSON(http.StatusInternalServerError, res)
}
