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

const (
	LoginActionGameLogin = "game_login"
)

type LoginRedirectData struct {
	Action             string
	GameLoginRequestID string
}

func (l LoginRedirectData) ToQuery() string {
	query := url.Values{}
	query.Set("action", l.Action)
	if l.GameLoginRequestID != "" {
		query.Set("game_login_request_id", l.GameLoginRequestID)
	}
	return query.Encode()
}

func (l LoginRedirectData) GetRedirectPath() string {
	switch l.Action {
	case LoginActionGameLogin:
		return "/game?" + url.Values{"id": []string{l.GameLoginRequestID}}.Encode()
	default:
		return "/feed"
	}
}

func ParseLoginRedirectData(redirect string) (LoginRedirectData, error) {
	query, err := url.ParseQuery(redirect)
	if err != nil {
		return LoginRedirectData{}, err
	}
	return LoginRedirectData{
		Action:             query.Get("action"),
		GameLoginRequestID: query.Get("game_login_request_id"),
	}, nil
}

func (c *LoginController) GetLogin(w http.ResponseWriter, r *http.Request) {
	redirect := r.URL.Query().Get("redirect")
	c.renderTemplate(w, &templates.LoginData{
		Redirect: redirect,
	})
}

func (c *LoginController) PostLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	redirect := r.FormValue("redirect")
	if username == "" || password == "" {
		c.renderTemplate(w, &templates.LoginData{
			Error:    "Username and password are required",
			Redirect: redirect,
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
			Error:    "Invalid username or password",
			Redirect: redirect,
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
	if redirect != "" {
		redirectData, err := ParseLoginRedirectData(redirect)
		if err != nil {
			http.Error(w, "Failed to parse redirect", http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, redirectData.GetRedirectPath(), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/feed", http.StatusSeeOther)
}
