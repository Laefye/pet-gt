package templates

import (
	"gt/internal/repository"
	"html/template"
)

type AuthenticatedData struct {
	User *repository.User
}

func parseAuthenticatedTemplate(files ...string) *template.Template {
	authenticated := []string{
		"web/templates/partial/nav.html",
		"web/templates/layout/authenticated.html",
	}
	all := make([]string, 0, len(baseTemplates)+len(authenticated)+len(files))
	all = append(all, baseTemplates...)
	all = append(all, authenticated...)
	all = append(all, files...)
	return template.Must(template.ParseFiles(all...))
}
