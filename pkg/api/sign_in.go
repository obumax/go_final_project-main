package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"go1f/pkg/api/middleware"
)

// структура входящего JSON
type signInRequest struct {
	Password string `json:"password"`
}

// signInHandler сверяет пароль из тела с секретом и выдает JWT-куки (/api/signin)
func signInHandler(w http.ResponseWriter, r *http.Request) {
	var req signInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("ошибка декодирования JSON: %v", err))
		return
	}

	secret := middleware.GetTokenSecret()
	// если аутентификация отключена (секрет пуст), просто OK без токена
	if secret == "" {
		writeJSON(w, http.StatusOK, map[string]any{"статус": "аутентификация отключена"})
		return
	}
	if req.Password != secret {
		writeError(w, http.StatusUnauthorized, "неверный пароль")
		return
	}

	// формируется токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "токен не подписан")
		return
	}

	// выставляются куки
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    signed,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})
	writeJSON(w, http.StatusOK, map[string]any{"статус": "ok"})
}
