package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "userID"
const RoleKey contextKey = "role"

type Claims struct {
	UserID int    `json:"userId"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")

	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(secret))
}

func Auth(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				http.Error(
					w,
					"Authorization header required",
					http.StatusUnauthorized,
				)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)

			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(
					w,
					"Invalid authorization format",
					http.StatusUnauthorized,
				)
				return
			}

			tokenString := parts[1]

			claims := &Claims{}

			token, err := jwt.ParseWithClaims(
				tokenString,
				claims,
				func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("JWT_SECRET")), nil
				},
			)

			if err != nil || !token.Valid {
				http.Error(
					w,
					"Invalid token",
					http.StatusUnauthorized,
				)
				return
			}

			if claims.Role != requiredRole {
				http.Error(
					w,
					"Forbidden",
					http.StatusForbidden,
				)
				return
			}

			ctx := context.WithValue(
				r.Context(),
				UserIDKey,
				claims.UserID,
			)

			ctx = context.WithValue(
				ctx,
				RoleKey,
				claims.Role,
			)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		})
	}
}
