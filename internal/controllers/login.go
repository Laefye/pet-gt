package controllers

import (
	"errors"
	"gt/internal/repository"
	"gt/internal/services"
	"gt/internal/templates"
	"net/http"
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
	c.renderTemplate(w, &templates.LoginData{})
}

func (c *LoginController) PostLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		c.renderTemplate(w, &templates.LoginData{Error: "Username and password are required"})
		return
	}
	session, err := c.authService.Login(r.Context(), services.LoginRequest{
		Username:  username,
		Password:  password,
		UserAgent: r.UserAgent(),
	})
	if err == nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    session.ID,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(24 * time.Hour),
		})
		http.Redirect(w, r, "/feed", http.StatusSeeOther)
		return
	}
	c.renderTemplate(w, &templates.LoginData{Error: "Invalid username or password"})
}

func (c *LoginController) Authenticate(r *http.Request) (*repository.User, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil, nil
		}
		return nil, err
	}
	session, err := c.authService.Authenticate(r.Context(), cookie.Value)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, nil
	}
	return session.User, nil
}
