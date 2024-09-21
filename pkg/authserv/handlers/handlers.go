package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SergeyIvanovDevelop/tss-tools/pkg/authserv/auth"
	"github.com/SergeyIvanovDevelop/tss-tools/pkg/authserv/repository"
	log "github.com/SergeyIvanovDevelop/tss-tools/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

var pkgLog log.Log

const pkgName string = "tss-tools/pkg/authserv/handlers"

type Credentials struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

type ValidateRequest struct {
	Token string `json:"token"`
}

type RevokeRequest struct {
	Token string `json:"token"`
}

func Register(repo repository.AuthRepository) http.HandlerFunc {
	fncLogger := log.AddLoggerFields(pkgLog, pkgName, log.Fields{
		"func": "Register",
	})
	return func(w http.ResponseWriter, r *http.Request) {
		fncLogger.Debug("Start")
		var creds Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			fncLogger.Error("Bad request:", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if creds.Password == "" || creds.Username == "" {
			fncLogger.Error("Empty username or password:", err)
			http.Error(w, "Empty username or password", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
		if err != nil {
			fncLogger.Error("Error process password:", err)
			http.Error(w, "Error process password", http.StatusInternalServerError)
			return
		}

		err = repo.CreateUser(creds.Username, string(hashedPassword))
		if err != nil {
			fncLogger.Error("Could not register user:", err)
			http.Error(w, "Could not register user", http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fncLogger.Debug("Finished")
	}
}

func Login(repo repository.AuthRepository) http.HandlerFunc {
	fncLogger := log.AddLoggerFields(pkgLog, pkgName, log.Fields{
		"func": "Login",
	})
	return func(w http.ResponseWriter, r *http.Request) {
		fncLogger.Debug("Start")
		var creds Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			fncLogger.Error("Bad request:", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if creds.Username == "" || creds.Password == "" {
			fncLogger.Error("Bad request:", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		storedPasswordHash, err := repo.GetUser(creds.Username)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(creds.Password)) != nil {
			fncLogger.Error("Unauthorized:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		accessToken, refreshToken, err := auth.GenerateToken(creds.Username)
		if err != nil {
			fncLogger.Error("Could not generate token:", err)
			http.Error(w, "Could not generate token", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
		if err != nil {
			fncLogger.Error("Error encoding json:", err)
			return
		}

		fncLogger.Debug("Finished")
	}
}

func Revoke(repo repository.AuthRepository) http.HandlerFunc {
	fncLogger := log.AddLoggerFields(pkgLog, pkgName, log.Fields{
		"func": "Revoke",
	})
	return func(w http.ResponseWriter, r *http.Request) {
		var revokeReq RevokeRequest
		err := json.NewDecoder(r.Body).Decode(&revokeReq)
		if err != nil || revokeReq.Token == "" {
			fncLogger.Error("Invalid request:", err)
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		claims, err := auth.ValidateToken(revokeReq.Token)
		if err != nil {
			fncLogger.Error("Invalid token:", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		err = repo.AddToBlacklist(revokeReq.Token, time.Unix(claims.ExpiresAt.Unix(), 0))
		if err != nil {
			fncLogger.Error("Failed to revoke token:", err)
			http.Error(w, "Failed to revoke token", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(map[string]string{"message": "Token successfully revoked"})
		if err != nil {
			fncLogger.Error("Error encoding json:", err)
			return
		}
	}
}

func Validate(repo repository.AuthRepository) http.HandlerFunc {
	fncLogger := log.AddLoggerFields(pkgLog, pkgName, log.Fields{
		"func": "Validate",
	})
	return func(w http.ResponseWriter, r *http.Request) {
		fncLogger.Debug("Start")
		var validateReq ValidateRequest
		err := json.NewDecoder(r.Body).Decode(&validateReq)
		if err != nil || validateReq.Token == "" {
			fncLogger.Error("Invalid request:", err)
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		if repo.IsInBlacklist(validateReq.Token) {
			fncLogger.Error("Token is revoked", err)
			http.Error(w, "Token is revoked", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateToken(validateReq.Token)
		if err != nil {
			fncLogger.Error("Invalid token:", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(map[string]string{
			"message": "Token is valid",
			"userID":  claims.Username,
		})
		if err != nil {
			fncLogger.Error("Error encoding json:", err)
			return
		}
		fncLogger.Debug("Finished")
	}
}
