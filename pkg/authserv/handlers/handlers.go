package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SergeyIvanovDevelop/tss-tools/pkg/authserv/auth"
	"github.com/SergeyIvanovDevelop/tss-tools/pkg/authserv/repository"

	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(repo repository.AuthRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		err = repo.CreateUser(creds.Username, creds.Password)
		if err != nil {
			http.Error(w, "Could not register user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func Login(repo repository.AuthRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		storedPasswordHash, err := repo.GetUser(creds.Username)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(creds.Password)) != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		accessToken, refreshToken, err := auth.GenerateToken(creds.Username)
		if err != nil {
			http.Error(w, "Could not generate token", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}

type RevokeRequest struct {
	Token string `json:"token"`
}

func Revoke(repo repository.AuthRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var revokeReq RevokeRequest
		err := json.NewDecoder(r.Body).Decode(&revokeReq)
		if err != nil || revokeReq.Token == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		claims, err := auth.ValidateToken(revokeReq.Token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		err = repo.AddToBlacklist(revokeReq.Token, time.Unix(claims.ExpiresAt.Unix(), 0))
		if err != nil {
			http.Error(w, "Failed to revoke token", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Token successfully revoked"})
	}
}

type ValidateRequest struct {
	Token string `json:"token"`
}

func Validate(repo repository.AuthRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var validateReq ValidateRequest
		err := json.NewDecoder(r.Body).Decode(&validateReq)
		if err != nil || validateReq.Token == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		if repo.IsInBlacklist(validateReq.Token) {
			http.Error(w, "Token is revoked", http.StatusUnauthorized)
			return
		}

		_, err = auth.ValidateToken(validateReq.Token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Token is valid"})
	}
}
