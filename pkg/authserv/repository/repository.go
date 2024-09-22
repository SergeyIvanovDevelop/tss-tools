package repository

import (
	"time"

	"github.com/SergeyIvanovDevelop/tss-tools/pkg/authserv/auth"
)

type AuthRepository interface {
	CreateUser(username, password string) error
	GetUser(username string) (string, error)
	AddToBlacklist(token string, expiration time.Time) error
	IsInBlacklist(token string) bool
	CleanExpiredTokens() error
	ValidateToken(tokenString string) (*auth.Claims, error)
}
