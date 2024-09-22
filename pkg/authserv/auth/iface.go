package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type Authenticator interface {
	GenerateToken(username string) (string, string, error)
	ValidateToken(tokenString string) (*Claims, error)
}
