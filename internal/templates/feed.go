package templates

type FeedData struct {
	AuthenticatedData
}

var FeedTemplate = parseAuthenticatedTemplate(
	"web/templates/page/feed.html",
)
