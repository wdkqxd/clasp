package helpers

import (
	"html"
	"regexp"
	"strings"
)

var reTag = regexp.MustCompile(`(?s)<[^>]+>`)

func CleanText(s string) string {
	s = reTag.ReplaceAllString(s, "")
	return strings.Join(strings.Fields(html.UnescapeString(s)), " ")
}
