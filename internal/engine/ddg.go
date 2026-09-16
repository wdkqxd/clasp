package engine

import (
	"context"
	"encoding/json"
	"errors"
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
	reVQD   = regexp.MustCompile(`vqd="([^"]+)"`)
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

func (d *DDG) ImageSearch(ctx context.Context, query string, limit, offset int) ([]ImageResult, error) {
	vqd, err := d.getVQD(ctx, query)
	if err != nil {
		return nil, err
	}

	u := fmt.Sprintf("https://duckduckgo.com/i.js?o=json&q=%s&vqd=%s", url.QueryEscape(query), url.QueryEscape(vqd))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", BrowserUserAgent)
	req.Header.Set("Referer", "https://duckduckgo.com/")
	req.Header.Set("Origin", "https://duckduckgo.com")

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ddg images request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ddg images returned %s", resp.Status)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, MaxBodyReadLimit))
	if err != nil {
		return nil, err
	}

	var data struct {
		Results []struct {
			Image      string `json:"image"`
			Title      string `json:"title"`
			Source     string `json:"source"`
			ThumbURL   string `json:"thumbnail"`
			Width      int    `json:"width"`
			Height     int    `json:"height"`
			SourceURL  string `json:"source_url"`
		} `json:"results"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("ddg images parse error: %w", err)
	}

	out := make([]ImageResult, 0, len(data.Results))
	for _, r := range data.Results {
		imgURL := r.Image
		if imgURL == "" {
			imgURL = r.SourceURL
		}
		if imgURL == "" {
			continue
		}
		out = append(out, ImageResult{
			Title:     r.Title,
			URL:       imgURL,
			Thumbnail: r.ThumbURL,
			Width:     r.Width,
			Height:    r.Height,
			Engine:    "ddg",
		})
		if len(out) >= limit {
			break
		}
	}
	if len(out) == 0 {
		return nil, errors.New("no image results found in ddg")
	}
	return out, nil
}

func (d *DDG) getVQD(ctx context.Context, query string) (string, error) {
	u, err := url.Parse("https://duckduckgo.com/")
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("q", query)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", BrowserUserAgent)

	resp, err := d.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, MaxBodyReadLimit))
	if err != nil {
		return "", err
	}

	m := reVQD.FindStringSubmatch(string(body))
	if len(m) < 2 {
		return "", errors.New("ddg: vqd token not found")
	}
	return m[1], nil
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
