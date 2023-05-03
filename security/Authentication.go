package security

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func Validate(c *gin.Context) {
	token, err := VerifyToken(c)

	if err != nil {
		// return err
	}

	// if
	claims, ok := token.Claims.(jwt.MapClaims)
	// !ok && !token.Valid {
	// 	return err
	// }
	logrus.Info("TOKEN")
	logrus.Info(claims)
	logrus.Info(token)
	logrus.Info(ok)

	// return nil
}

func VerifyToken(c *gin.Context) (*jwt.Token, error) {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("SECRET"), nil
	})

	if err != nil {
		return nil, err
	}

	logrus.Info("==============")
	logrus.Info(token)
	return token, nil
}

func ExtractToken(c *gin.Context) string {
	// keys := c.Request.URL.Query()
	// token := keys.Get("token")

	// if token != "" {
	// 	return token
	// }

	c.Writer.Header()
	authHeader := c.GetHeader("Authorization")
	bearerToken := strings.Split(authHeader, " ")

	if len(bearerToken) == 2 {
		logrus.Info(bearerToken[1])
		return bearerToken[1]
	} else {
		return ""
	}
}
