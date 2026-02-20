package templates

import "gt/internal/repository"

type GameData struct {
	GameLoginRequest *repository.GameLoginRequest
	User             *repository.User
	Error            string
}

var GameTemplate = parseTemplate(
	"web/templates/page/game.html",
)
