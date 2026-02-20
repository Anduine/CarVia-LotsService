package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
)

func AuthMiddleware(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			slog.Debug("Заголовок авторизації відсутній")
			http.Error(w, "Не авторизовано", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := UserIDFromToken(tokenString)
		if err != nil {
			slog.Debug("Помилка авторизації", "err", err.Error())
			http.Error(w, "Не авторизовано", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		next(w, r.WithContext(ctx))
	})
}

func OptionalAuthMiddleware(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := 0
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			if id, err := UserIDFromToken(tokenString); err != nil {
				// slog.Debug("У OptionalAuthMiddleware нема авторизації", "err", err.Error())
				userID = 0
			} else {
				userID = id
			}
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		next(w, r.WithContext(ctx))
	})
}

func AuthMiddlewareHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			slog.Debug("Заголовок авторизації відсутній")
			http.Error(w, "Не авторизовано", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := UserIDFromToken(tokenString)
		if err != nil {
			slog.Debug("Помилка авторизації", "err", err.Error())
			http.Error(w, "Не авторизовано", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func OptionalAuthMiddlewareHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := 0
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			if id, err := UserIDFromToken(tokenString); err != nil {
				// slog.Debug("У OptionalAuthMiddleware нема авторизації", "err", err.Error())
				userID = 0
			} else {
				userID = id
			}
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
