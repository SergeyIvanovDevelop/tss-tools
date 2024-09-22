package repository

import "github.com/SergeyIvanovDevelop/tss-tools/pkg/authserv/auth"

func ValidateToken(tokenString string) (*auth.Claims, error) {
	return auth.ValidateToken(tokenString)
}
