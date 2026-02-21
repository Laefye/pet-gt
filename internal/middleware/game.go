package middleware

import (
	"context"
	"gt/internal/repository"
	"gt/internal/services"
	"net/http"
)

const gameLoginContextKey contextKey = "gamelogin"

func GameLoginFromContext(ctx context.Context) *repository.GameLogin {
	gameLogin, _ := ctx.Value(gameLoginContextKey).(*repository.GameLogin)
	return gameLogin
}

func RequireGameLogin(gameService *services.GameService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Game-Login-ID")
		token := r.Header.Get("X-Game-Login-Token")
		if id == "" || token == "" {
			http.Error(w, "Missing game login credentials", http.StatusUnauthorized)
			return
		}
		gameLogin, err := gameService.AuthenticateGameLogin(r.Context(), id, token)
		if err != nil {
			http.Error(w, "Invalid game login credentials", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), gameLoginContextKey, gameLogin)
		next(w, r.WithContext(ctx))
	}
}
