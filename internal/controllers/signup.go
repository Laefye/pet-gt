package controllers

import (
	"errors"
	"gt/internal/services"
	"gt/internal/templates"
	"net/http"
)

type SignupController struct {
	authService *services.AuthService
}

func NewSignupController(authService *services.AuthService) *SignupController {
	return &SignupController{authService: authService}
}

func (c *SignupController) renderTemplate(w http.ResponseWriter, data *templates.SignupData) {
	err := templates.SignupTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *SignupController) GetSignup(w http.ResponseWriter, r *http.Request) {
	c.renderTemplate(w, &templates.SignupData{})
}

func (c *SignupController) PostSignup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		c.renderTemplate(w, &templates.SignupData{Error: "Username and password are required"})
		return
	}
	_, err := c.authService.Signup(r.Context(), services.SignupRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		var signupErr *services.SignupError
		if errors.As(err, &signupErr) {
			c.renderTemplate(w, &templates.SignupData{Error: signupErr.Message})
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
