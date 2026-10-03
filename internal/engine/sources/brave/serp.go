package brave

import (
	"errors"
	"regexp"
	"strings"

	"clasp/internal/engine"
	"clasp/internal/engine/helpers"
)

const serpMarker = `{type:"data",data:`

var (
	reSSRTitle   = regexp.MustCompile(`(?s)class="title search-snippet-title[^"]*"[^>]*title="([^"]+)"`)
	reSSRSnippet = regexp.MustCompile(`(?s)class="generic-snippet[^"]*">.*?class="content[^"]*">(.*?)</div>`)
	reAnchorHref = regexp.MustCompile(`href="([^"]+)"`)
)

func parseSerpData(doc string) (map[string]any, error) {
	for start := 0; ; {
		i := strings.Index(doc[start:], serpMarker)
		if i < 0 {
			break
		}
		i += start + len(serpMarker)
		end, ok := helpers.MatchObject(doc, i)
		if !ok {
			start = i
			continue
		}
		v, err := helpers.ParseJSLiteral(doc[i:end])
		start = i
		if err != nil {
			continue
		}
		obj, ok := v.(map[string]any)
		if !ok {
			continue
		}
		if _, has := obj["body"]; has {
			return obj, nil
		}
	}
	return nil, errors.New("nothing was found. most likely, the page layout has changed")
}

func webResults(data map[string]any, limit int) ([]engine.Result, error) {
	body, _ := data["body"].(map[string]any)
	response, _ := body["response"].(map[string]any)
	web, _ := response["web"].(map[string]any)
	raw, _ := web["results"].([]any)

	out := make([]engine.Result, 0, len(raw))
	for _, item := range raw {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		title, _ := m["title"].(string)
		resultURL, _ := m["url"].(string)
		if title == "" && resultURL == "" {
			continue
		}
		snippet, _ := m["description"].(string)
		out = append(out, engine.Result{
			Title:     title,
			URL:       resultURL,
			Snippet:   snippet,
			Engine:    "brave",
			Sitelinks: parseSitelinks(m["cluster"]),
		})
		if len(out) >= limit {
			break
		}
	}
	if len(out) == 0 {
		return nil, errors.New("no web results found in brave payload")
	}
	return out, nil
}

const maxSitelinks = 4

func parseSitelinks(v any) []engine.Sitelink {
	raw, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]engine.Sitelink, 0, len(raw))
	for _, item := range raw {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		linkURL, _ := m["url"].(string)
		if !strings.HasPrefix(linkURL, "http://") && !strings.HasPrefix(linkURL, "https://") {
			continue
		}
		title, _ := m["title"].(string)
		out = append(out, engine.Sitelink{Title: title, URL: linkURL})
		if len(out) >= maxSitelinks {
			break
		}
	}
	return out
}

func ssrResults(doc string, limit int) ([]engine.Result, error) {
	idxs := reSSRTitle.FindAllStringIndex(doc, -1)
	if len(idxs) == 0 {
		return nil, errors.New("no brave results found in page")
	}

	out := make([]engine.Result, 0, len(idxs))
	for _, m := range idxs {
		title := reSSRTitle.FindStringSubmatch(doc[m[0]:m[1]])[1]

		open := strings.LastIndex(doc[:m[0]], "<a")
		if open < 0 {
			continue
		}
		cl := strings.Index(doc[open:], ">")
		if cl < 0 {
			continue
		}
		anchor := doc[open : open+cl]
		// skip sponsored cards
		if strings.Contains(doc[open:m[1]], "ad-info-label") || strings.Contains(doc[open:m[1]], "Sponsored") {
			continue
		}

		href := reAnchorHref.FindStringSubmatch(anchor)
		if len(href) < 2 || href[1] == "" {
			continue
		}
		resultURL := href[1]
		if strings.HasPrefix(resultURL, "//") {
			resultURL = "https:" + resultURL
		}
		if !strings.HasPrefix(resultURL, "http://") && !strings.HasPrefix(resultURL, "https://") {
			continue
		}

		r := engine.Result{Title: helpers.CleanText(title), URL: resultURL, Engine: "brave"}
		if s := reSSRSnippet.FindStringSubmatch(doc[m[1]:]); len(s) > 1 {
			r.Snippet = helpers.CleanText(s[1])
		}
		out = append(out, r)
		if len(out) >= limit {
			break
		}
	}
	if len(out) == 0 {
		return nil, errors.New("no web results found in brave page")
	}
	return out, nil
}
