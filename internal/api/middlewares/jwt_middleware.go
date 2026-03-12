package middlewares

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"rest-api/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("Bearer")

		if err != nil {
			http.Error(w, "Authorization Header Missing", http.StatusUnauthorized)
			return
		}

		jwtSecret := os.Getenv("JWT_SECRET")

		parsedToken, err := jwt.Parse(token.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(jwtSecret), nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				http.Error(w, "Token Expired", http.StatusUnauthorized)
				return
			}

			if errors.Is(err, jwt.ErrTokenMalformed) {
				http.Error(w, "Token Malformed", http.StatusUnauthorized)
				return
			}

			utils.ErrorHandler(err, "")
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if parsedToken.Valid {
			log.Println("Token Valid")
		} else {
			http.Error(w, "Invalid Login Token", http.StatusUnauthorized)
			log.Println("Invalid Token:", token.Value)
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid Login Token", http.StatusUnauthorized)
			log.Println("Invalid Token:", token.Value)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKey("role"), claims["role"])
		ctx = context.WithValue(ctx, ContextKey("expiresAt"), claims["exp"])
		ctx = context.WithValue(r.Context(), ContextKey("username"), claims["user"])
		ctx = context.WithValue(r.Context(), ContextKey("userId"), claims["uid"])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
