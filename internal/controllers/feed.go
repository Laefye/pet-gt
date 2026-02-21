package controllers

import (
	"gt/internal/middleware"
	"gt/internal/repository"
	"gt/internal/services"
	"gt/internal/templates"
	"net/http"
)

type FeedController struct {
	achievementService *services.AchievementService
}

func NewFeedController(achievementService *services.AchievementService) *FeedController {
	return &FeedController{achievementService: achievementService}
}

func (c *FeedController) GetFeed(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	data := templates.FeedData{
		AuthenticatedData: templates.AuthenticatedData{
			User: user,
		},
	}
	achievements, err := c.achievementService.GetAchievementsByUserID(r.Context(), user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, achievement := range achievements {
		data.Achievements = append(data.Achievements, templates.AchievementData{
			Name:      repository.AchievementName(achievement.Name).String(),
			ImageURL:  repository.AchievementName(achievement.Name).ImageURL(),
			CreatedAt: achievement.CreatedAt,
		})
	}
	err = templates.FeedTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
