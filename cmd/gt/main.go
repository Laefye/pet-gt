package main

import (
	"fmt"
	"gt/internal/controllers"
	"gt/internal/middleware"
	"gt/internal/repository"
	"gt/internal/services"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "12345678"),
		getEnv("DB_NAME", "gt"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_SSLMODE", "disable"),
		getEnv("DB_TIMEZONE", "Europe/Moscow"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	if err := db.AutoMigrate(&repository.User{}, &repository.Session{}, &repository.GameLogin{}, &repository.GameLoginRequest{}); err != nil {
		log.Fatal("failed to migrate database: ", err)
	}

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	gameLoginRepo := repository.NewGameLoginRepository(db)
	gameLoginRequestRepo := repository.NewGameLoginRequestRepository(db)

	authService := services.NewAuthService(userRepo, sessionRepo)
	gameService := services.NewGameService(gameLoginRepo, gameLoginRequestRepo)

	signupCtrl := controllers.NewSignupController(authService)
	loginCtrl := controllers.NewLoginController(authService)
	feedCtrl := controllers.NewFeedController()
	gameCtrl := controllers.NewGameController(gameService)
	profileCtrl := controllers.NewProfileController(authService)

	auth := func(next http.HandlerFunc) http.HandlerFunc {
		return middleware.RequireAuth(authService, next)
	}
	optAuth := func(next http.HandlerFunc) http.HandlerFunc {
		return middleware.OptionalAuth(authService, next)
	}
	noAuth := func(next http.HandlerFunc) http.HandlerFunc {
		return middleware.NoAuth(authService, "/login", next)
	}

	mux := http.NewServeMux()
	mux.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("web/public"))))

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/feed", http.StatusSeeOther)
	})

	mux.HandleFunc("GET /signup", noAuth(signupCtrl.GetSignup))
	mux.HandleFunc("POST /signup", noAuth(signupCtrl.PostSignup))
	mux.HandleFunc("GET /login", loginCtrl.GetLogin)
	mux.HandleFunc("POST /login", loginCtrl.PostLogin)

	mux.HandleFunc("GET /feed", auth(feedCtrl.GetFeed))

	mux.HandleFunc("POST /api/game", gameCtrl.CreateGameLoginRequest)
	mux.HandleFunc("GET /api/game", gameCtrl.GetGameLoginState)
	mux.HandleFunc("GET /game", optAuth(gameCtrl.GetGameLoginPage))
	mux.HandleFunc("POST /game", auth(gameCtrl.PostGameLogin))
	mux.HandleFunc("GET /profile/logout", auth(profileCtrl.Logout))

	addr := getEnv("LISTEN_ADDR", "localhost:8080")
	log.Printf("server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
