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

func BadRequest(c *gin.Context) {
	res := BadRequestError{
		Success: false,
		Message: "Bad Request",
		Data:    nil,
	}

	c.JSON(http.StatusBadRequest, res)
}
