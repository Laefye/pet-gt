package templates

import (
	"gt/internal/repository"
	"html/template"
)

var baseLoginedTemplates = append(
	baseTemplates,
	"web/templates/partial/nav.html",
	"web/templates/layout/logined.html",
)

type LoginedData struct {
	User *repository.User
}

func parseLoginedTemplate(files ...string) *template.Template {
	allFiles := append(baseLoginedTemplates, files...)
	return template.Must(template.ParseFiles(allFiles...))
}
