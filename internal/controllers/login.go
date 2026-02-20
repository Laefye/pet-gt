package controllers

import (
	"gt/internal/services"
	"gt/internal/templates"
	"net/http"
	"net/url"
	"time"
)

type LoginController struct {
	authService *services.AuthService
}

func NewLoginController(authService *services.AuthService) *LoginController {
	return &LoginController{authService: authService}
}

func (c *LoginController) renderTemplate(w http.ResponseWriter, data *templates.LoginData) {
	err := templates.LoginTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *LoginController) GetLogin(w http.ResponseWriter, r *http.Request) {
	gameLoginRequestID := r.URL.Query().Get("game_login_request_id")
	c.renderTemplate(w, &templates.LoginData{
		GameLoginRequestID: gameLoginRequestID,
	})
}

func (c *LoginController) PostLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	gameLoginRequestID := r.FormValue("game_login_request_id")
	if username == "" || password == "" {
		c.renderTemplate(w, &templates.LoginData{
			Error:              "Username and password are required",
			GameLoginRequestID: gameLoginRequestID,
		})
		return
	}
	session, err := c.authService.Login(r.Context(), services.LoginRequest{
		Username:  username,
		Password:  password,
		UserAgent: r.UserAgent(),
	})
	if err != nil {
		c.renderTemplate(w, &templates.LoginData{
			Error:              "Invalid username or password",
			GameLoginRequestID: gameLoginRequestID,
		})
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	})
	if gameLoginRequestID != "" {
		query := url.Values{}
		query.Set("id", gameLoginRequestID)
		http.Redirect(w, r, "/game?"+query.Encode(), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/feed", http.StatusSeeOther)
}
