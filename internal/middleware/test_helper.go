package middleware

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateTestToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	})
	return token.SignedString([]byte("api-workbench-jwt-secret-2026"))
}
