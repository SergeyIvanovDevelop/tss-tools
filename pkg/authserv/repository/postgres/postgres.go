package postgres

import (
	"context"
	"time"

	"github.com/SergeyIvanovDevelop/tss-tools/pkg/authserv/auth"
	"github.com/jackc/pgx/v4"
)

type PostgresAuthRepository struct {
	conn *pgx.Conn
}

func NewPostgresRepository(connString string) (*PostgresAuthRepository, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	return &PostgresAuthRepository{conn: conn}, nil
}

func (repo *PostgresAuthRepository) Close() {
	repo.conn.Close(context.Background())
}

func (repo *PostgresAuthRepository) CreateUser(username, hashedPassword string) error {
	_, err := repo.conn.Exec(context.Background(),
		"INSERT INTO users_auth (username, password) VALUES ($1, $2)", username, hashedPassword)
	return err
}

func (repo *PostgresAuthRepository) GetUser(username string) (string, error) {
	var passwordHash string
	err := repo.conn.QueryRow(context.Background(),
		"SELECT password FROM users_auth WHERE username=$1", username).Scan(&passwordHash)
	if err != nil {
		return "", err
	}
	return passwordHash, nil
}

func (repo *PostgresAuthRepository) AddToBlacklist(token string, expiration time.Time) error {
	_, err := repo.conn.Exec(context.Background(),
		"INSERT INTO token_blacklist (token, expires_at) VALUES ($1, $2)", token, expiration)
	return err
}

func (repo *PostgresAuthRepository) IsInBlacklist(token string) bool {
	var exists bool
	err := repo.conn.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM token_blacklist WHERE token=$1)", token).Scan(&exists)
	return err == nil && exists
}

func (repo *PostgresAuthRepository) CleanExpiredTokens() error {
	_, err := repo.conn.Exec(context.Background(),
		"DELETE FROM token_blacklist WHERE expires_at < NOW()")
	return err
}

func (repo *PostgresAuthRepository) ValidateToken(tokenString string) (*auth.Claims, error) {
	return auth.ValidateToken(tokenString)
}
