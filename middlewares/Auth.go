package middlewares

import (
	"motor/exceptions"
	"motor/security"

	"github.com/gin-gonic/gin"
)

func SetupAuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := security.GetAuthentication(c)

		if err != nil {
			exceptions.AuthorizeException(c, err.Error())
			c.Abort()
			return
		}

		c.Next()
	}
}
