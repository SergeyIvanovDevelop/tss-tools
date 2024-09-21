package middleware

import (
	"net/http"
	"strings"

	"github.com/SergeyIvanovDevelop/tss-tools/pkg/authserv/auth"
	log "github.com/SergeyIvanovDevelop/tss-tools/pkg/logger"
)

var pkgLog log.Log

const pkgName string = "tss-tools/pkg/authserv/middleware"

func JWTAuthentication(next http.Handler) http.Handler {
	fncLogger := log.AddLoggerFields(pkgLog, pkgName, log.Fields{
		"func": "JWTAuthentication",
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			fncLogger.Error("No header 'Authorization'")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		_, err := auth.ValidateToken(token)
		if err != nil {
			fncLogger.Errorf("Not valid token '%s'", token)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
