package middleware

import (
	"context"
	"errors"
	"gt/internal/repository"
	"gt/internal/services"
	"net/http"
)

type contextKey string

const userContextKey contextKey = "user"

func UserFromContext(ctx context.Context) *repository.User {
	user, _ := ctx.Value(userContextKey).(*repository.User)
	return user
}

func RequireAuth(authService *services.AuthService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := authenticateRequest(r, authService)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next(w, r.WithContext(ctx))
	}
}

func OptionalAuth(authService *services.AuthService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := authenticateRequest(r, authService)
		if user != nil {
			ctx := context.WithValue(r.Context(), userContextKey, user)
			r = r.WithContext(ctx)
		}
		next(w, r)
	}
}

func authenticateRequest(r *http.Request, authService *services.AuthService) *repository.User {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil
		}
		return nil
	}
	session, err := authService.Authenticate(r.Context(), cookie.Value)
	if err != nil || session == nil {
		return nil
	}
	return session.User
}
