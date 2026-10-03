package brave

import (
	"html"
	"net/url"
	"regexp"
	"strings"

	"clasp/internal/engine"
	"clasp/internal/engine/helpers"
)

var (
	reInfoboxAnchor = regexp.MustCompile(`(?s)<a[^>]+href="(https://[^"]+\.wikipedia\.org/wiki/[^"]+)"[^>]*>(.*?)</a>`)
	reInfoboxURL    = regexp.MustCompile(`"url":"(https://[^"]+\.wikipedia\.org/wiki/[^"]+)"`)
)

func parseInfobox(doc string) *engine.Infobox {
	if i := strings.Index(doc, `id="infobox"`); i >= 0 {
		sec := doc[i:]
		if end := strings.Index(sec, `</section>`); end >= 0 {
			sec = sec[:end]
		}
		if m := reInfoboxAnchor.FindStringSubmatch(sec); len(m) == 3 {
			title := helpers.CleanText(m[2])
			if title == "" {
				title = wikiTitleFromURL(m[1])
			}
			return &engine.Infobox{Title: title, URL: m[1]}
		}
	}
	if i := strings.Index(doc, `infobox:{type:"graph"`); i >= 0 {
		window := doc[i:]
		if len(window) > 8192 {
			window = window[:8192]
		}
		if m := reInfoboxURL.FindStringSubmatch(window); len(m) == 2 {
			return &engine.Infobox{Title: wikiTitleFromURL(m[1]), URL: m[1]}
		}
	}
	return nil
}

func wikiTitleFromURL(raw string) string {
	i := strings.LastIndex(raw, "/wiki/")
	if i < 0 {
		return ""
	}
	title, err := url.PathUnescape(raw[i+len("/wiki/"):])
	if err != nil {
		title = raw[i+len("/wiki/"):]
	}
	return strings.ReplaceAll(title, "_", " ")
}

var (
	reDiscussionTitle   = regexp.MustCompile(`(?s)summary title="([^"]*)"`)
	reDiscussionSnippet = regexp.MustCompile(`(?s)<div class="content[^>]*>(.*?)<a class="source`)
)

func parseDiscussions(doc string) []engine.Discussion {
	i := strings.Index(doc, `id="discussions"`)
	if i < 0 {
		return nil
	}
	sec := doc[i:]
	if end := strings.Index(sec, `</section>`); end >= 0 {
		sec = sec[:end]
	}

	seen := make(map[string]struct{})
	out := make([]engine.Discussion, 0, 8)
	rest := sec
	for {
		start := strings.Index(rest, `<div class="discussions-item`)
		if start < 0 {
			break
		}
		rest = rest[start:]
		marker := `<div class="discussions-item`
		next := strings.Index(rest[len(marker):], marker)
		end := len(rest)
		if next >= 0 {
			end = next + len(marker)
		}
		item := rest[:end]
		rest = rest[end:]

		t := reDiscussionTitle.FindStringSubmatch(item)
		if len(t) < 2 {
			continue
		}
		title := strings.TrimSpace(helpers.CleanText(t[1]))
		if title == "" {
			continue
		}
		url := firstHTTPLink(item)
		if url == "" {
			continue
		}
		if _, dup := seen[url]; dup {
			continue
		}
		seen[url] = struct{}{}
		snippet := ""
		if s := reDiscussionSnippet.FindStringSubmatch(item); len(s) > 1 {
			snippet = strings.TrimSpace(helpers.CleanText(s[1]))
		}
		out = append(out, engine.Discussion{Title: title, URL: url, Snippet: snippet})
		if len(out) >= 8 {
			break
		}
	}
	return out
}

func firstHTTPLink(s string) string {
	for _, m := range reAnchorHref.FindAllStringSubmatch(s, -1) {
		if len(m) < 2 {
			continue
		}
		href := html.UnescapeString(m[1])
		if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
			return href
		}
	}
	return ""
}
