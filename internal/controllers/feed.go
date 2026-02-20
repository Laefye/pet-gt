package controllers

import (
	"gt/internal/repository"
	"gt/internal/templates"
	"net/http"
)

type FeedController struct {
}

func NewFeedController() *FeedController {
	return &FeedController{}
}

func (c *FeedController) GetFeed(user *repository.User, w http.ResponseWriter, r *http.Request) {
	data := templates.FeedData{
		LoginedData: templates.LoginedData{
			User: user,
		},
	}
	err := templates.FeedTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
