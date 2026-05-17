package sanitizer

import "github.com/microcosm-cc/bluemonday"

func SanitizeContent(html string) string {
	var policy = bluemonday.UGCPolicy()
	policy.AllowAttrs("href").OnElements("a")
	policy.RequireNoFollowOnLinks(true)

	return policy.Sanitize(html)
}
