package controllers

import (
	"gt/internal/middleware"
	"gt/internal/templates"
	"net/http"
)

type FeedController struct{}

func NewFeedController() *FeedController {
	return &FeedController{}
}

func (c *FeedController) GetFeed(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	data := templates.FeedData{
		AuthenticatedData: templates.AuthenticatedData{
			User: user,
		},
	}
	err := templates.FeedTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
