package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lumenratio/finalTODO/internal/logger"
)

type PassString struct {
	Password string `json:"password"`
}

type Claims struct {
	Hash string
	jwt.RegisteredClaims
}

func passHash(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

func authGenTokenHandler(w http.ResponseWriter, r *http.Request) {
	// смотрим наличие пароля
	pass := os.Getenv("TODO_PASSWORD")
	if len(pass) > 0 {
		secret := PassString{}
		// Decode JSON body
		err := json.NewDecoder(r.Body).Decode(&secret)
		if err != nil {
			writeError(w, "malformed auth request", http.StatusBadRequest)
			return
		}
		if i := strings.Compare(pass, secret.Password); i != 0 {
			writeError(w, "wrong password", http.StatusUnauthorized)
			return
		}
		// Создаем JWT
		payload := Claims{
			Hash: passHash(secret.Password),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
		// Подписываем
		sigToken, err := token.SignedString([]byte(secret.Password))
		if err != nil {
			writeError(w, "failed sign jwt", http.StatusInternalServerError)
			return
		}
		logger.Info.Println("new token:", JWTData{Token: sigToken})
		writeJson(w, JWTData{Token: sigToken})
	}
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var ctoken string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				ctoken = cookie.Value
			}
			var valid bool
			// здесь код для валидации и проверки JWT-токена
			payload := Claims{}
			token, err := jwt.ParseWithClaims(ctoken, &payload, func(t *jwt.Token) (interface{}, error) {
				return []byte(pass), nil
			})
			if err != nil {
				writeError(w, "token corrupted", http.StatusBadRequest)
				return
			}
			// Валидируем и сверям хеш
			if token.Valid {
				valid = true
				if i := strings.Compare(payload.Hash, passHash(pass)); i != 0 {
					writeError(w, "password chaged, token not valid anymore", http.StatusUnauthorized)
					return
				}
			}
			if !valid {
				// возвращаем ошибку авторизации 401
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
