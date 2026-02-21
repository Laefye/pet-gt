package controllers

import (
	"gt/internal/middleware"
	"gt/internal/services"
	"net/http"
)

type ProfileController struct {
	authService *services.AuthService
}

func NewProfileController(authService *services.AuthService) *ProfileController {
	return &ProfileController{authService: authService}
}

func (c *ProfileController) Logout(w http.ResponseWriter, r *http.Request) {
	session := middleware.SessionFromContext(r.Context())
	err := c.authService.Logout(r.Context(), session.ID)
	if err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
