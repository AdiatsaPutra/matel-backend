package payloads

import "github.com/gin-gonic/gin"

type Respons struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func HandleSuccess(c *gin.Context, data interface{}, message string, status int) {
	res := Respons{
		Success: true,
		Message: message,
		Data:    data,
	}

	c.JSON(status, res)
}
