package mware

import (
	"context"
	"gitlab.com/ictisagora/backend/internal/lib/jwt"
	"net/http"
	"strings"
)

// Определим типы ключей для контекста:
type contextKey string

const (
	ContextUserUID   contextKey = "userUID"
	ContextUserEmail contextKey = "userEmail"
)

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				http.Error(w, "No token provided", http.StatusUnauthorized)
				return
			}

			uid, email, err := jwt.ParseToken(token, jwtSecret)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Кладём uid и email в context
			ctx := context.WithValue(r.Context(), ContextUserUID, uid)
			ctx = context.WithValue(ctx, ContextUserEmail, email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
