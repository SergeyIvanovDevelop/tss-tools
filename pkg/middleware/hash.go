package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	log "github.com/SergeyIvanovDevelop/tss-tools/pkg/logger"
)

const hashHeaderName = "HashSHA256"

func getHash(key string, payload []byte) string {
	var fncLogger = log.AddLoggerFields(pkgLog, pkgName, log.Fields{
		"func": "getHash",
	})
	if key == "" {
		fncLogger.Warn("'key' is empty")
		return ""
	}
	if len(payload) == 0 {
		fncLogger.Warn("'payload' is empty")
		return ""
	}
	h := hmac.New(sha256.New, []byte(key))
	h.Write(payload)
	hashSum := h.Sum(nil)
	hexStrHashSum := hex.EncodeToString(hashSum[:])
	return hexStrHashSum
}

// responseRecorder перехватывает запись в ResponseWriter
type responseRecorder struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
	rr.body.Write(b)
	return rr.ResponseWriter.Write(b)
}

func buildHashMiddleware(key string) func(next http.Handler) http.Handler {
	var fncLogger = log.AddLoggerFields(pkgLog, pkgName, log.Fields{
		"func": "buildHashMiddleware (return func)",
	})
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Проверка hash у POST-запросов
			if r.Method == http.MethodPost {
				bodyBytes, err := io.ReadAll(r.Body)
				if err != nil {
					fncLogger.Error("Failed to read request body")
					http.Error(w, "Failed to read request body", http.StatusInternalServerError)
					return
				}
				defer r.Body.Close()

				// Восстанавливаем тело запроса для дальнейшего использования
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				hashString := getHash(key, bodyBytes)

				clientHash := r.Header.Get(hashHeaderName)
				// Игнорируем проверку, если пустой заголовок 'hashHeaderName' пришел
				if clientHash != "" {
					fncLogger.Warn("Заголовок ХЕША - НЕ пустой!")
					if clientHash != hashString {
						fncLogger.Errorf("Hash mismatch: clientHash='%s' | hashString='%s'", clientHash, hashString)
						http.Error(w, "Hash mismatch", http.StatusBadRequest)
						return
					}
				}

			}

			rr := &responseRecorder{ResponseWriter: w, body: &bytes.Buffer{}}

			next.ServeHTTP(rr, r)

			hashString := getHash(key, rr.body.Bytes())
			w.Header().Set(hashHeaderName, hashString)

		})
	}
}
