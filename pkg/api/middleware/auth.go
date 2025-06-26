package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

var tokenSecret string

// SetTokenSecret задает секрет для подписи JWT
// Если секрет пустой, пропускаются все запросы
func SetTokenSecret(secret string) {
	tokenSecret = secret
}

func GetTokenSecret() string {
	return tokenSecret
}

type errResp struct {
	Error string `json:"error"`
}

// Auth проверяет JWT-куки "token"
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if tokenSecret != "" {
			cookie, err := r.Cookie("token")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(errResp{"требуется аутентификация"})
				return
			}
			token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrTokenUnverifiable
				}
				return []byte(tokenSecret), nil
			})
			if err != nil || !token.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(errResp{"требуется аутентификация"})
				return
			}
		}
		next(w, r)
	}
}
