package main

import (
	"gt/internal/controllers"
	"gt/internal/repository"
	"gt/internal/services"
	"gt/internal/templates"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func renderSignupError(w http.ResponseWriter, message string) {
	data := templates.SignupData{Error: message}
	err := templates.SignupTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	dsn := "host=localhost user=postgres password=12345678 dbname=gt port=5432 sslmode=disable TimeZone=Europe/Moscow"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&repository.User{}, &repository.Session{})
	if err != nil {
		panic(err)
	}

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	authService := services.NewAuthService(userRepo, sessionRepo)

	signupController := controllers.NewSignupController(authService)
	loginController := controllers.NewLoginController(authService)
	feedController := controllers.NewFeedController()

	mux := http.NewServeMux()
	mux.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("web/public"))))
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
	mux.HandleFunc("GET /signup", func(w http.ResponseWriter, r *http.Request) {
		signupController.GetSignup(w, r)
	})
	mux.HandleFunc("POST /signup", func(w http.ResponseWriter, r *http.Request) {
		signupController.PostSignup(w, r)
	})
	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		loginController.GetLogin(w, r)
	})
	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		loginController.PostLogin(w, r)
	})
	mux.HandleFunc("GET /feed", func(w http.ResponseWriter, r *http.Request) {
		user, err := loginController.Authenticate(r)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		feedController.GetFeed(user, w, r)
	})

	http.ListenAndServe("localhost:8080", mux)
}
