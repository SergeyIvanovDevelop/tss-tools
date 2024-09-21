package repository

import (
	"time"
)

type AuthRepository interface {
	CreateUser(username, password string) error
	GetUser(username string) (string, error)
	AddToBlacklist(token string, expiration time.Time) error
	IsInBlacklist(token string) bool
	CleanExpiredTokens() error
}
