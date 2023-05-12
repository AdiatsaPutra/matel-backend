package security

import (
	"fmt"
	"motor/exceptions"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)
var jwtKey = []byte("SECRET")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the JWT token from the request header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			exceptions.Unauthorized(c, "Missing authorization header")
			c.Abort()
			return
		}
		
		tokenString := authHeader[len("Bearer "):]
		
		// Parse the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			
			return jwtKey, nil
		})
		
		// Verify the JWT token
		if err != nil || !token.Valid {
			exceptions.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}
		
		// Pass the username from the token to the next handlers
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID := uint(claims["user_id"].(float64))
			c.Set("user_id", userID)
			} else {
			exceptions.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}

		c.Next()
	}
}
