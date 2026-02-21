package templates

import "time"

type AchievementData struct {
	Name      string
	ImageURL  string
	CreatedAt time.Time
}

type FeedData struct {
	AuthenticatedData
	Achievements []AchievementData
}

var FeedTemplate = parseAuthenticatedTemplate(
	"web/templates/partial/achievement.html",
	"web/templates/page/feed.html",
)
