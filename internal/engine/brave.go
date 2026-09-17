package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const (
	braveWebEndpoint = "https://search.brave.com/search"
	braveWebMarker   = `{type:"data",data:`
)

type Brave struct {
	client HTTPClient

	mu      sync.Mutex
	infobox *Infobox
}

func NewBrave(client HTTPClient) *Brave {
	return &Brave{client: client}
}

func (b *Brave) Name() string { return "brave" }

func (b *Brave) Infobox() *Infobox {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.infobox
}

func (b *Brave) ImageSearch(ctx context.Context, query string, limit, offset int) ([]ImageResult, error) {
	u, err := url.Parse("https://search.brave.com/images")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("q", query)
	q.Set("source", "web")
	q.Set("spellcheck", "0")
	q.Set("search_lang", "en")
	q.Set("country", "ALL")
	if offset > 0 {
		q.Set("offset", strconv.Itoa(offset))
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	setBraveHeaders(req)

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("brave images returned %s", resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, MaxBodyReadLimit))
	if err != nil {
		return nil, err
	}
	doc := string(body)

	data, err := parseSerpData(doc)
	if err != nil {
		return nil, err
	}
	return braveImageResults(data, limit)
}

func braveImageResults(data map[string]any, limit int) ([]ImageResult, error) {
	body, _ := data["body"].(map[string]any)
	response, _ := body["response"].(map[string]any)
	raw, _ := response["results"].([]any)

	out := make([]ImageResult, 0, len(raw))
	for _, item := range raw {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		title, _ := m["title"].(string)
		sourceURL, _ := m["url"].(string)

		// thumbnail contains the small image
		thumb := ""
		thumbW, thumbH := 0, 0
		if thumbObj, ok := m["thumbnail"].(map[string]any); ok {
			thumb, _ = thumbObj["src"].(string)
			thumbW = intNum(thumbObj["width"])
			thumbH = intNum(thumbObj["height"])
		}

		// properties contains the full-size image
		imgURL := sourceURL
		width, height := thumbW, thumbH
		if props, ok := m["properties"].(map[string]any); ok {
			if u, ok := props["url"].(string); ok && u != "" {
				imgURL = u
			}
			width = intNum(props["width"])
			height = intNum(props["height"])
		}

		if imgURL == "" && thumb == "" {
			continue
		}

		out = append(out, ImageResult{
			Title:     title,
			URL:       sourceURL,
			SourceURL: imgURL,
			Thumbnail: thumb,
			Width:     width,
			Height:    height,
			Engine:    "brave",
		})
		if len(out) >= limit {
			break
		}
	}
	if len(out) == 0 {
		return nil, errors.New("no image results found in brave payload")
	}
	return out, nil
}

func intNum(v any) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case int64:
		return int(n)
	case json.Number:
		i, _ := n.Int64()
		return int(i)
	}
	return 0
}

func (b *Brave) Search(ctx context.Context, query string, limit, offset int) ([]Result, error) {
	b.mu.Lock()
	b.infobox = nil
	b.mu.Unlock()

	u, err := url.Parse(braveWebEndpoint)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("q", query)
	q.Set("source", "web")
	q.Set("spellcheck", "0")
	q.Set("search_lang", "en")
	q.Set("country", "ALL")
	if offset > 0 {
		q.Set("offset", strconv.Itoa(offset))
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	setBraveHeaders(req)

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("brave web returned %s", resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, MaxBodyReadLimit))
	if err != nil {
		return nil, err
	}
	doc := string(body)

	data, serpErr := parseSerpData(doc)
	b.mu.Lock()
	b.infobox = parseBraveInfobox(doc)
	b.mu.Unlock()

	if serpErr == nil {
		return webResults(data, limit)
	}
	return parseSSRResults(doc, limit)
}

func setBraveHeaders(req *http.Request) {
	req.Header.Set("User-Agent", BrowserUserAgent)
	req.Header.Set("Cookie", "theme=dark; country=all; useLocation=0")
	req.Header.Set("sec-ch-ua", `"Chromium";v="133", "Not(A:Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("sec-gpc", "1")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("referer", "https://search.brave.com/")
}

func parseSerpData(doc string) (map[string]any, error) {
	for start := 0; ; {
		i := strings.Index(doc[start:], braveWebMarker)
		if i < 0 {
			break
		}
		i += start + len(braveWebMarker)
		end, ok := matchObject(doc, i)
		if !ok {
			start = i
			continue
		}
		v, err := parseJSLiteral(doc[i:end])
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
	return nil, errors.New("brave serp data not found: layout changed or blocked")
}

// returns the offset past the literal at s[i], skipping nested brackets/strings/escapes
func matchObject(s string, i int) (int, bool) {
	if i >= len(s) || (s[i] != '{' && s[i] != '[') {
		return 0, false
	}
	open := s[i]
	close := byte('}')
	if open == '[' {
		close = ']'
	}
	depth := 0
	inStr := false
	var quote byte
	for ; i < len(s); i++ {
		c := s[i]
		if inStr {
			if c == '\\' {
				i++
			} else if c == quote {
				inStr = false
			}
			continue
		}
		switch c {
		case '"', '\'':
			inStr = true
			quote = c
		case open:
			depth++
		case close:
			depth--
			if depth == 0 {
				return i + 1, true
			}
		}
	}
	return 0, false
}

func webResults(data map[string]any, limit int) ([]Result, error) {
	body, _ := data["body"].(map[string]any)
	response, _ := body["response"].(map[string]any)
	web, _ := response["web"].(map[string]any)
	raw, _ := web["results"].([]any)

	out := make([]Result, 0, len(raw))
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
		out = append(out, Result{
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

func parseSitelinks(v any) []Sitelink {
	raw, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]Sitelink, 0, len(raw))
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
		out = append(out, Sitelink{Title: title, URL: linkURL})
		if len(out) >= maxSitelinks {
			break
		}
	}
	return out
}

var (
	reBraveSSRTitle   = regexp.MustCompile(`(?s)class="title search-snippet-title[^"]*"[^>]*title="([^"]+)"`)
	reBraveSSRSnippet = regexp.MustCompile(`(?s)class="generic-snippet[^"]*">.*?class="content[^"]*">(.*?)</div>`)
	reAnchorHref      = regexp.MustCompile(`href="([^"]+)"`)
)

func parseSSRResults(doc string, limit int) ([]Result, error) {
	idxs := reBraveSSRTitle.FindAllStringIndex(doc, -1)
	if len(idxs) == 0 {
		return nil, errors.New("no brave results found in page")
	}

	out := make([]Result, 0, len(idxs))
	for _, m := range idxs {
		title := reBraveSSRTitle.FindStringSubmatch(doc[m[0]:m[1]])[1]

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

		r := Result{Title: cleanText(title), URL: resultURL, Engine: "brave"}
		if s := reBraveSSRSnippet.FindStringSubmatch(doc[m[1]:]); len(s) > 1 {
			r.Snippet = cleanText(s[1])
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

var (
	reBraveInfoboxAnchor = regexp.MustCompile(`(?s)<a[^>]+href="(https://[^"]+\.wikipedia\.org/wiki/[^"]+)"[^>]*>(.*?)</a>`)
	reBraveInfoboxURL    = regexp.MustCompile(`"url":"(https://[^"]+\.wikipedia\.org/wiki/[^"]+)"`)
)

func parseBraveInfobox(doc string) *Infobox {
	if i := strings.Index(doc, `id="infobox"`); i >= 0 {
		sec := doc[i:]
		if end := strings.Index(sec, `</section>`); end >= 0 {
			sec = sec[:end]
		}
		if m := reBraveInfoboxAnchor.FindStringSubmatch(sec); len(m) == 3 {
			title := cleanText(m[2])
			if title == "" {
				title = wikiTitleFromURL(m[1])
			}
			return &Infobox{Title: title, URL: m[1]}
		}
	}
	if i := strings.Index(doc, `infobox:{type:"graph"`); i >= 0 {
		window := doc[i:]
		if len(window) > 8192 {
			window = window[:8192]
		}
		if m := reBraveInfoboxURL.FindStringSubmatch(window); len(m) == 2 {
			return &Infobox{Title: wikiTitleFromURL(m[1]), URL: m[1]}
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
