package templates

type FeedData struct {
	LoginedData
}

var FeedTemplate = parseLoginedTemplate(
	"web/templates/page/feed.html",
)
