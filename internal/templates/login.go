package templates

type LoginData struct {
	Error    string
	Redirect string
}

var LoginTemplate = parseTemplate(
	"web/templates/page/login.html",
)

type SignupData struct {
	Error string
}

var SignupTemplate = parseTemplate(
	"web/templates/page/signup.html",
)
