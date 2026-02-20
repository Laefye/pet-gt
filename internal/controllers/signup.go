package controllers

import (
	"errors"
	"gt/internal/services"
	"gt/internal/templates"
	"net/http"
)

type SignupController struct {
	userService *services.AuthService
}

func NewSignupController(userService *services.AuthService) *SignupController {
	return &SignupController{userService: userService}
}

func (c *SignupController) renderTemplate(w http.ResponseWriter, data *templates.SignupData) {
	err := templates.SignupTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *SignupController) GetSignup(w http.ResponseWriter, r *http.Request) {
	c.renderTemplate(w, nil)
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
	_, err := c.userService.Signup(r.Context(), services.SignupRequest{
		Username: username,
		Password: password,
	})
	if err == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	var signupError *services.SignupError
	if errors.As(err, &signupError) {
		c.renderTemplate(w, &templates.SignupData{Error: signupError.Message})
		return
	}
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}
