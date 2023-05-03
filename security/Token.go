package security

import (
	// "time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userId uint) (string, error) {
	// timer := time.Now().Add(AuthTokenValidTime).Unix()
	secret := "SECRET"
	claims := jwt.MapClaims{}

	claims["authorized"] = true
	claims["user_id"] = userId
	// claims["exp"] = timer

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return tokenString, err
}
