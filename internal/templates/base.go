package templates

import "html/template"

var baseTemplates = []string{
	"web/templates/layout/base.html",
}

func parseTemplate(files ...string) *template.Template {
	all := make([]string, 0, len(baseTemplates)+len(files))
	all = append(all, baseTemplates...)
	all = append(all, files...)
	return template.Must(template.ParseFiles(all...))
}
