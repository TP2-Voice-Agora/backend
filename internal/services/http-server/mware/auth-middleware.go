package mware

import (
	"context"
	"github.com/TP2-Voice-Agora/backend/internal/lib/jwt"
	i "github.com/TP2-Voice-Agora/backend/internal/services/interfaces"
	"log/slog"
	"net/http"
	"strings"
)

type contextKey string

const (
	ContextUserUID   contextKey = "userUID"
	ContextUserEmail contextKey = "userEmail"
)

func AuthMiddleware(jwtSecret string, log *slog.Logger, s i.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				log.Error("No authorization header found" + authHeader)
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

				log.Error("Failed to parse token", slog.String("error", err.Error()), slog.String("token", token))
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			u, err := s.GetUserByUID(uid)
			if err != nil {
				log.Error("Failed to get user by uid", slog.String("error", err.Error()))
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}

			if u.ReAuth == true {
				log.Warn("User re-auth, token will be reset", slog.String("uid", uid), slog.String("email", email))
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
