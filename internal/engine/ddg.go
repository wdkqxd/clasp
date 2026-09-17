package engine

import (
	"context"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const ddgEndpoint = "https://html.duckduckgo.com/html/"

var (
	reTitle = regexp.MustCompile(`(?i)<a[^>]*class="[^"]*result__a[^"]*"[^>]*href="([^"]+)"[^>]*>(.*?)</a>`)
	reSnip  = regexp.MustCompile(`(?i)<a[^>]*class="[^"]*result__snippet[^"]*"[^>]*>(.*?)</a>`)
	reTag   = regexp.MustCompile(`(?s)<[^>]+>`)
)

type DDG struct {
	client HTTPClient
}

func NewDDG(client HTTPClient) *DDG {
	return &DDG{client: client}
}

func (d *DDG) Name() string { return "ddg" }

func (d *DDG) Search(ctx context.Context, query string, limit, offset int) ([]Result, error) {
	u, err := url.Parse(ddgEndpoint)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("q", query)
	if offset > 0 {
		q.Set("s", strconv.Itoa(offset*30))
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", BrowserUserAgent)

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("ddg blocked the request: %s", resp.Status)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ddg returned %s", resp.Status)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, MaxBodyReadLimit))
	if err != nil {
		return nil, err
	}
	doc := string(body)

	titles := reTitle.FindAllStringSubmatch(doc, -1)
	snips := reSnip.FindAllStringSubmatch(doc, -1)

	results := make([]Result, 0, len(titles))
	for i, m := range titles {
		if len(m) < 3 {
			continue
		}
		target, err := resolveDDGURL(m[1])
		if err != nil || isDDGAd(target) {
			continue
		}
		r := Result{
			Title:  cleanText(m[2]),
			URL:    target,
			Engine: d.Name(),
		}
		if i < len(snips) && len(snips[i]) > 1 {
			r.Snippet = cleanText(snips[i][1])
		}
		results = append(results, r)
		if len(results) >= limit {
			break
		}
	}
	return results, nil
}

func resolveDDGURL(href string) (string, error) {
	if strings.Contains(href, "uddg=") {
		u, err := url.Parse(href)
		if err == nil {
			if target := u.Query().Get("uddg"); target != "" {
				return target, nil
			}
		}
	}
	switch {
	case strings.HasPrefix(href, "//"):
		return "https:" + href, nil
	case strings.HasPrefix(href, "http://"), strings.HasPrefix(href, "https://"):
		return href, nil
	}
	return "", fmt.Errorf("unsupported ddg link: %q", href)
}

func cleanText(s string) string {
	s = reTag.ReplaceAllString(s, "")
	return strings.Join(strings.Fields(html.UnescapeString(s)), " ")
}

// sponsored-result redirects (y.js links)
func isDDGAd(target string) bool {
	return strings.Contains(target, "ad_domain=") || strings.Contains(target, "/y.js")
}
