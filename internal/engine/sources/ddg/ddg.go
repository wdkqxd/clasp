package ddg

import (
	"context"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"clasp/internal/engine"
	"clasp/internal/engine/helpers"
)

const (
	endpoint = "https://html.duckduckgo.com/html/"
	pageSize = 30
	// ddg only walks a few pages per query, no need to remember more cursors
	maxCursors = 32
)

var (
	reTitle = regexp.MustCompile(`(?i)<a[^>]*class="[^"]*result__a[^"]*"[^>]*href="([^"]+)"[^>]*>(.*?)</a>`)
	reSnip  = regexp.MustCompile(`(?i)<a[^>]*class="[^"]*result__snippet[^"]*"[^>]*>(.*?)</a>`)

	reNextForm = regexp.MustCompile(`(?s)<form[^>]*action="/html/"[^>]*>.*?</form>`)
	reField    = regexp.MustCompile(`(?is)<input[^>]*>`)
	reFieldKey = regexp.MustCompile(`(?i)name="([^"]*)"`)
	reFieldVal = regexp.MustCompile(`(?i)value="([^"]*)"`)
	reVqd      = regexp.MustCompile(`vqd="(4-[0-9-]+)"`)
)

type DDG struct {
	client engine.HTTPClient

	mu      sync.Mutex
	cursors map[string]url.Values // "Next" form per query, ddg's only way to paginate
	vqds    map[string]string     // vqd tokens, keyed by query
}

func New(client engine.HTTPClient) *DDG {
	return &DDG{
		client:  client,
		cursors: make(map[string]url.Values),
		vqds:    make(map[string]string),
	}
}

func (d *DDG) Name() string { return "ddg" }

func (d *DDG) Search(ctx context.Context, query string, limit, page int) (engine.SearchResult, error) {
	doc, err := d.fetchPage(ctx, query, page)
	if err != nil {
		return engine.SearchResult{}, err
	}
	d.remember(query, doc)
	return engine.SearchResult{Results: parseResults(doc, limit)}, nil
}

func parseResults(doc string, limit int) []engine.Result {
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
			Engine: "ddg",
		}
		if i < len(snips) && len(snips[i]) > 1 {
			r.Snippet = helpers.CleanText(snips[i][1])
		}
		results = append(results, r)
		if len(results) >= limit {
			break
		}
	}
	return results
}

func (d *DDG) fetchPage(ctx context.Context, query string, page int) (string, error) {
	if page == 0 {
		return d.get(ctx, url.Values{"q": {query}})
	}
	if form := d.getCursor(query); form != nil {
		if doc, err := d.post(ctx, form); err == nil {
			return doc, nil
		}
	}
	form := url.Values{"q": {query}, "s": {strconv.Itoa(page * pageSize)}, "v": {"l"}, "o": {"json"}}
	if vqd := d.getVQD(query); vqd != "" {
		form.Set("vqd", vqd)
	}
	return d.get(ctx, form)
}

func (d *DDG) get(ctx context.Context, form url.Values) (string, error) {
	return d.do(ctx, http.MethodGet, endpoint+"?"+form.Encode(), nil)
}

func (d *DDG) post(ctx context.Context, form url.Values) (string, error) {
	return d.do(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
}

func (d *DDG) do(ctx context.Context, method, rawURL string, body io.Reader) (string, error) {
	const maxRedirects = 2
	for redirect := 0; ; redirect++ {
		var reader io.Reader
		if redirect == 0 {
			reader = body
		}
		req, err := http.NewRequestWithContext(ctx, method, rawURL, reader)
		if err != nil {
			return "", err
		}
		req.Header.Set("User-Agent", engine.BrowserUserAgent)
		req.Header.Set("Referer", endpoint)
		if method == http.MethodPost {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}

		resp, err := d.client.Do(req)
		if err != nil {
			return "", err
		}
		doc, err := readDoc(resp)
		next := resp.Header.Get("Location")
		if err != nil || redirect >= maxRedirects || next == "" {
			return doc, err
		}
		loc, err := resp.Location()
		if err != nil {
			return doc, nil
		}
		method, rawURL, body = http.MethodGet, loc.String(), nil
	}
}

func readDoc(resp *http.Response) (string, error) {
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests {
		return "", fmt.Errorf("blocked the request: %s", resp.Status)
	}
	if resp.StatusCode < 300 || resp.StatusCode > 399 {
		if resp.StatusCode != http.StatusOK {
			return "", errors.New(resp.Status)
		}
		body, err := io.ReadAll(io.LimitReader(resp.Body, engine.MaxBodyReadLimit))
		if err != nil {
			return "", err
		}
		return string(body), nil
	}
	io.Copy(io.Discard, resp.Body)
	return "", nil
}

func (d *DDG) remember(query, doc string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if m := reVqd.FindStringSubmatch(doc); m != nil {
		d.vqds[query] = m[1]
	}
	form := nextForm(doc)
	if form == nil {
		return
	}
	form.Set("q", query)
	if vqd := d.vqds[query]; vqd != "" {
		form.Set("vqd", vqd)
	}
	if len(d.cursors) >= maxCursors {
		clear(d.cursors)
	}
	d.cursors[query] = form
}

func (d *DDG) getCursor(query string) url.Values {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.cursors[query]
}

func (d *DDG) getVQD(query string) string {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.vqds[query]
}

func nextForm(doc string) url.Values {
	forms := reNextForm.FindAllString(doc, -1)
	if len(forms) == 0 {
		return nil
	}
	form := url.Values{}
	for _, field := range reField.FindAllString(forms[len(forms)-1], -1) {
		name := reFieldKey.FindStringSubmatch(field)
		if name == nil || name[1] == "" {
			continue
		}
		value := ""
		if v := reFieldVal.FindStringSubmatch(field); v != nil {
			value = html.UnescapeString(v[1])
		}
		form.Set(name[1], value)
	}
	if len(form) == 0 {
		return nil
	}
	return form
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
