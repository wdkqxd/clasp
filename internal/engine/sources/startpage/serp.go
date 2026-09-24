package startpage

import (
	"errors"
	"html"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"metasearch/internal/engine"
)

var (
	reSpTitle       = regexp.MustCompile(`(?s)<a class="result-title result-link[^"]*" href="([^"]+)"[^>]*>\s*<h2 class="wgl-title[^"]*"[^>]*>(.*?)</h2>\s*</a>`)
	reSpDescription = regexp.MustCompile(`(?s)<p class="description[^"]*"[^>]*>(.*?)</p>`)
	reSpTags        = regexp.MustCompile(`(?s)<[^>]+>`)
	reSpPageForm    = regexp.MustCompile(`(?s)<form action="/sp/search" method="post"[^>]*aria-label="go to page ([0-9]+)"[^>]*>(.*?)</form>`)
	reSpInputTag    = regexp.MustCompile(`(?s)<input\b[^>]*>`)
	reSpAttr        = regexp.MustCompile(`([\w.\-]+)=("([^"]*)"|'([^']*)'|[^\s>]+)`)
)

func cleanText(raw string) string {
	raw = reSpTags.ReplaceAllString(raw, "")
	return html.UnescapeString(strings.TrimSpace(raw))
}

func parseResults(doc string, limit int) ([]engine.Result, error) {
	idxs := reSpTitle.FindAllStringSubmatchIndex(doc, -1)
	if len(idxs) == 0 {
		return nil, errors.New("no startpage results found in page")
	}

	out := make([]engine.Result, 0, len(idxs))
	for i, m := range idxs {
		resultURL := doc[m[2]:m[3]]
		if !strings.HasPrefix(resultURL, "http://") && !strings.HasPrefix(resultURL, "https://") {
			continue
		}
		r := engine.Result{
			Title:  cleanText(doc[m[4]:m[5]]),
			URL:    resultURL,
			Engine: "startpage",
		}

		segStart := m[1]
		segEnd := len(doc)
		if i+1 < len(idxs) {
			segEnd = idxs[i+1][0]
		}
		if d := reSpDescription.FindStringSubmatch(doc[segStart:segEnd]); len(d) > 1 {
			r.Snippet = cleanText(d[1])
		}

		out = append(out, r)
		if len(out) >= limit {
			break
		}
	}
	if len(out) == 0 {
		return nil, errors.New("no usable startpage results found")
	}
	return out, nil
}

func pageForm(doc string, page int) (url.Values, bool) {
	for _, m := range reSpPageForm.FindAllStringSubmatch(doc, -1) {
		p, err := strconv.Atoi(m[1])
		if err != nil || p != page {
			continue
		}
		form := url.Values{}
		for _, tag := range reSpInputTag.FindAllString(m[2], -1) {
			attrs := make(map[string]string)
			for _, a := range reSpAttr.FindAllStringSubmatch(tag, -1) {
				v := a[3]
				if v == "" {
					v = a[4]
				}
				attrs[a[1]] = v
			}
			if name, ok := attrs["name"]; ok {
				form.Set(name, attrs["value"])
			}
		}
		return form, true
	}
	return nil, false
}
