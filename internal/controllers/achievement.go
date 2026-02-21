package controllers

import (
	"encoding/json"
	"errors"
	"gt/internal/middleware"
	"gt/internal/repository"
	"gt/internal/services"
	"net/http"
)

type AchievementController struct {
	achievementService *services.AchievementService
}

type achievementResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	UserID string `json:"user_id"`
}

type achievementErrorResponse struct {
	Message string `json:"message"`
}

func NewAchievementController(achievementService *services.AchievementService) *AchievementController {
	return &AchievementController{achievementService: achievementService}
}

func (c *AchievementController) jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (c *AchievementController) AddAchievement(w http.ResponseWriter, r *http.Request) {
	name := repository.AchievementName(r.FormValue("name"))
	user := middleware.GameLoginFromContext(r.Context()).User
	if name == "" {
		c.jsonResponse(w, achievementErrorResponse{Message: "Name are required"}, http.StatusBadRequest)
		return
	}
	if !name.IsValid() {
		c.jsonResponse(w, achievementErrorResponse{Message: "Invalid achievement name"}, http.StatusBadRequest)
		return
	}
	achievement, err := c.achievementService.CreateAchievement(r.Context(), &repository.CreateAchievementRequest{
		Name:   name,
		UserID: user.ID,
	})
	if errors.Is(err, services.ErrAchievementAlreadyExists) {
		c.jsonResponse(w, achievementErrorResponse{Message: "Achievement already exists for user"}, http.StatusConflict)
		return
	} else if err != nil {
		c.jsonResponse(w, achievementErrorResponse{Message: "Failed to add achievement"}, http.StatusInternalServerError)
		return
	}
	c.jsonResponse(w, achievementResponse{
		ID:     achievement.ID,
		Name:   achievement.Name,
		UserID: achievement.UserID,
	}, http.StatusCreated)
}
