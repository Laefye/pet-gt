package templates

import "gt/internal/repository"

type GameData struct {
	GameLoginRequest *repository.GameLoginRequest
	User             *repository.User
	Error            string
}

var GameLoginTemplate = parseTemplate(
	"web/templates/page/game/login.html",
)
