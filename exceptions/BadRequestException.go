package exceptions

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BadRequestError struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func BadRequest(c *gin.Context, message string) {
	res := BadRequestError{
		Success: false,
		Message: message,
		Data:    nil,
	}

	c.JSON(http.StatusBadRequest, res)
}
