package security

import (
	"fmt"
	"motor/payloads"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ExtractAuthToken(c *gin.Context) (*payloads.UserDetail, error) {
	token, err := VerifyToken(c)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		username := fmt.Sprintf("%s", claims["username"])
		authorized := fmt.Sprintf("%s", claims["authorized"])

		if err != nil {
			return nil, err
		}

		authDetail := payloads.UserDetail{
			Username:  username,
			Authorize: authorized,
		}

		return &authDetail, nil
	}

	return nil, err
}
