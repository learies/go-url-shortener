package middlewares

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"
	"github.com/learies/go-url-shortener/internal/contextutils"
	"github.com/learies/go-url-shortener/internal/logger"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

func createUserID() string {
	userID := uuid.New().String()
	logger.Log.Info("Created new user ID", "userID", userID)
	return userID
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID string
		var tokenString string

		// Чтение токена из куки
		cookie, err := r.Cookie("token")
		if err == nil {
			tokenString = cookie.Value
		}

		// Если токен не передан нужно создать userID и создать для него токен
		if tokenString == "" {
			logger.Log.Error("No token provided")
			userID = createUserID()

			// Время жизни токена 1 минута
			expirationTime := time.Now().Add(1 * time.Minute)

			claims := &Claims{
				UserID: userID,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(expirationTime),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			tokenString, err := token.SignedString([]byte("qwerty"))
			if err != nil {
				http.Error(w, "Could not create token", http.StatusInternalServerError)
				return
			}

			// Устанавливаем токен в куки
			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    tokenString,
				Expires:  expirationTime,
				HttpOnly: true,
				Path:     "/",
			})
		} else {
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte("qwerty"), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			userID = claims.UserID
			logger.Log.Info("Got user ID from token in cookie", "userID", userID)
		}

		ctx := contextutils.WithUserID(r.Context(), userID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
