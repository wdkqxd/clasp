package ddg

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"metasearch/internal/engine"
	"metasearch/internal/engine/helpers"
)

const endpoint = "https://html.duckduckgo.com/html/"

var (
	reTitle = regexp.MustCompile(`(?i)<a[^>]*class="[^"]*result__a[^"]*"[^>]*href="([^"]+)"[^>]*>(.*?)</a>`)
	reSnip  = regexp.MustCompile(`(?i)<a[^>]*class="[^"]*result__snippet[^"]*"[^>]*>(.*?)</a>`)
)

type DDG struct {
	client engine.HTTPClient
}

func New(client engine.HTTPClient) *DDG {
	return &DDG{client: client}
}

func (d *DDG) Name() string { return "ddg" }

func (d *DDG) Search(ctx context.Context, query string, limit, offset int) ([]engine.Result, error) {
	u, err := url.Parse(endpoint)
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
	req.Header.Set("User-Agent", engine.BrowserUserAgent)

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

	body, err := io.ReadAll(io.LimitReader(resp.Body, engine.MaxBodyReadLimit))
	if err != nil {
		return nil, err
	}
	doc := string(body)

	titles := reTitle.FindAllStringSubmatch(doc, -1)
	snips := reSnip.FindAllStringSubmatch(doc, -1)

	results := make([]engine.Result, 0, len(titles))
	for i, m := range titles {
		if len(m) < 3 {
			continue
		}
		target, err := resolveURL(m[1])
		if err != nil || isAd(target) {
			continue
		}
		r := engine.Result{
			Title:  helpers.CleanText(m[2]),
			URL:    target,
			Engine: d.Name(),
		}
		if i < len(snips) && len(snips[i]) > 1 {
			r.Snippet = helpers.CleanText(snips[i][1])
		}
		results = append(results, r)
		if len(results) >= limit {
			break
		}
	}
	return results, nil
}

func resolveURL(href string) (string, error) {
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

// sponsored-result redirects (y.js links)
func isAd(target string) bool {
	return strings.Contains(target, "ad_domain=") || strings.Contains(target, "/y.js")
}
