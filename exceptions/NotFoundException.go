package exceptions

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NotFoundException(c *gin.Context, message string) {
	res := Error{
		Success: false,
		Message: message,
		Data: nil,
	}

	c.JSON(http.StatusNotFound, res)
}