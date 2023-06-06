package exceptions

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UnauthorizedExceptionError struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Unauthorized(c *gin.Context, message string) {
	res := UnauthorizedExceptionError{
		Success: false,
		Message: message,
		Data:    nil,
	}

	c.JSON(http.StatusUnauthorized, res)
}
