package ddg

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

	"clasp/internal/engine"
)

const imageEndpoint = "https://duckduckgo.com/i.js"

var reVQD = regexp.MustCompile(`vqd="([^"]+)"`)

type DDGImageResult struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	Image     string `json:"image"`
	Thumbnail string `json:"thumbnail"`
	Source    string `json:"source"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

type ddgImagePayload struct {
	Results []DDGImageResult `json:"results"`
}

func (d *DDG) ImageSearch(ctx context.Context, query string, limit, offset int) ([]engine.ImageResult, error) {
	vqd, err := d.vqd(ctx, query)
	if err != nil {
		return nil, err
	}

	args := url.Values{}
	args.Set("o", "json")
	args.Set("q", query)
	args.Set("l", "us-en")
	args.Set("vqd", vqd)
	args.Set("p", "1")
	args.Set("ct", "AT")
	if offset > 0 {
		args.Set("s", strconv.Itoa(offset*100))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageEndpoint+"?"+args.Encode(), nil)
	if err != nil {
		return nil, err
	}
	setImageHeaders(req)

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
		return nil, fmt.Errorf("blocked the request: %s", resp.Status)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	var payload ddgImagePayload
	if err := json.NewDecoder(io.LimitReader(resp.Body, engine.MaxBodyReadLimit)).Decode(&payload); err != nil {
		return nil, err
	}

	out := make([]engine.ImageResult, 0, len(payload.Results))
	for _, r := range payload.Results {
		if r.Thumbnail == "" || r.Image == "" {
			continue
		}
		// the raw image URL is hotlink-protected, serve it through DDG's proxy
		proxied := "https://external-content.duckduckgo.com/iu/?u=" + url.QueryEscape(r.Image)
		out = append(out, engine.ImageResult{
			Title:     r.Title,
			URL:       r.URL,
			SourceURL: proxied,
			Thumbnail: proxied,
			Width:     r.Width,
			Height:    r.Height,
			Engine:    d.Name(),
		})
		if len(out) >= limit {
			break
		}
	}
	if len(out) == 0 {
		return nil, errors.New("no image results found in ddg payload")
	}
	return out, nil
}

func (d *DDG) vqd(ctx context.Context, query string) (string, error) {
	d.mu.Lock()
	got := d.vqds[query]
	d.mu.Unlock()
	if got != "" {
		return got, nil
	}

	token, err := d.fetchVQD(ctx, query)
	if err != nil {
		return "", err
	}

	d.mu.Lock()
	d.vqds[query] = token
	d.mu.Unlock()
	return token, nil
}

func (d *DDG) fetchVQD(ctx context.Context, query string) (string, error) {
	u := "https://duckduckgo.com/?q=" + url.QueryEscape(query)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", engine.BrowserUserAgent)

	resp, err := d.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s while fetching vqd", resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", err
	}

	m := reVQD.FindSubmatch(body)
	if len(m) < 2 || len(m[1]) == 0 {
		return "", errors.New("no vqd value found in duckduckgo response")
	}
	return string(m[1]), nil
}

func setImageHeaders(req *http.Request) {
	req.Header.Set("User-Agent", engine.BrowserUserAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Referer", "https://duckduckgo.com/")
	req.Header.Set("Sec-GPC", "1")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Priority", "u=4")
}
