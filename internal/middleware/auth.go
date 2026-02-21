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
const sessionContextKey contextKey = "session"

func UserFromContext(ctx context.Context) *repository.User {
	user, _ := ctx.Value(userContextKey).(*repository.User)
	return user
}

func SessionFromContext(ctx context.Context) *repository.Session {
	session, _ := ctx.Value(sessionContextKey).(*repository.Session)
	return session
}

func RequireAuth(authService *services.AuthService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := authenticateRequest(r, authService)
		if session == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), userContextKey, session.User)
		ctx = context.WithValue(ctx, sessionContextKey, session)
		next(w, r.WithContext(ctx))
	}
}

func OptionalAuth(authService *services.AuthService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := authenticateRequest(r, authService)
		if session != nil {
			ctx := context.WithValue(r.Context(), userContextKey, session.User)
			ctx = context.WithValue(ctx, sessionContextKey, session)
			r = r.WithContext(ctx)
		}
		next(w, r)
	}
}

func NoAuth(authService *services.AuthService, redirectURL string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := authenticateRequest(r, authService)
		if session != nil {
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func authenticateRequest(r *http.Request, authService *services.AuthService) *repository.Session {
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
	return session
}
