package templates

import "html/template"

var baseTemplates = []string{
	"web/templates/layout/base.html",
}

func parseTemplate(files ...string) *template.Template {
	allFiles := append(baseTemplates, files...)
	return template.Must(template.ParseFiles(allFiles...))
}
