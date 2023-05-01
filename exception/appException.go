package exception

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func AppException(c *gin.Context, message string) {
	res := Error{
		Success: false,
		Message: message,
		Data:    nil,
	}

	c.JSON(http.StatusInternalServerError, res)
}
