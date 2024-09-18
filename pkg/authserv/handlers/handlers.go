package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/SergeyIvanovDevelop/tss-tools/pkg/authserv/auth"
	"github.com/SergeyIvanovDevelop/tss-tools/pkg/authserv/repository"

	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

func Register(repo repository.AuthRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("[Register] START")
		var creds Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			fmt.Println("[Register] Bad request: ", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if creds.Password == "" || creds.Username == "" {
			fmt.Println("[Register] Empty username or password: ", err)
			http.Error(w, "Empty username or password", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("[Register] Error process password: ", err)
			http.Error(w, "Error process password", http.StatusInternalServerError)
			return
		}

		err = repo.CreateUser(creds.Username, string(hashedPassword))
		if err != nil {
			fmt.Println("[Register] Could not register user: ", err)
			http.Error(w, "Could not register user", http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Println("[Register] OK")
	}
}

func Login(repo repository.AuthRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("[Login] START")
		var creds Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			fmt.Println("[Login] Bad request: ", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		storedPasswordHash, err := repo.GetUser(creds.Username)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(creds.Password)) != nil {
			fmt.Println("[Login] Unauthorized: ", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		accessToken, refreshToken, err := auth.GenerateToken(creds.Username)
		if err != nil {
			fmt.Println("[Login] Could not generate token: ", err)
			http.Error(w, "Could not generate token", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
		fmt.Println("[Login] OK")
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
		fmt.Println("[Validate] START")
		var validateReq ValidateRequest
		err := json.NewDecoder(r.Body).Decode(&validateReq)
		if err != nil || validateReq.Token == "" {
			fmt.Println("[Validate] Invalid request")
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		if repo.IsInBlacklist(validateReq.Token) {
			fmt.Println("[Validate] Token is revoked")
			http.Error(w, "Token is revoked", http.StatusUnauthorized)
			return
		}

		_, err = auth.ValidateToken(validateReq.Token)
		if err != nil {
			fmt.Println("[Validate] Invalid token")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Token is valid"})
		fmt.Println("[Validate] OK")
	}
}
