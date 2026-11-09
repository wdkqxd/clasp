package wikipedia

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode"

	"clasp/internal/engine"
	"clasp/internal/engine/helpers"
)

type Article struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Extract     string `json:"extract,omitempty"`
	URL         string `json:"url"`
	Thumbnail   string `json:"thumbnail,omitempty"`
}

const (
	perLanguageTimeout = 3 * time.Second
	cacheTTL           = 5 * time.Minute
	negativeCacheTTL   = 90 * time.Second
	maxCacheEntries    = 512
	upstreamInterval   = 150 * time.Millisecond
	throttleCooldown   = 60 * time.Second
	maxRedirects       = 2
	userAgent          = "clasp/1.0 (https://github.com/wdkq/clasp)"
	summaryURL         = "https://%s.wikipedia.org/api/rest_v1/page/summary/%s"
)

type cacheEntry struct {
	article *Article
	expires time.Time
}

type Client struct {
	http engine.HTTPClient

	mu           sync.Mutex
	cache        map[string]cacheEntry
	nextRequest  time.Time
	blockedUntil time.Time
}

func NewClient(httpClient engine.HTTPClient) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 5 * time.Second}
	}
	return &Client{http: httpClient, cache: make(map[string]cacheEntry)}
}

var errThrottled = errors.New("wikipedia: throttled")

func (c *Client) waitSlot(ctx context.Context) error {
	c.mu.Lock()
	now := time.Now()
	if c.nextRequest.Before(now) {
		c.nextRequest = now
	}
	at := c.nextRequest
	c.nextRequest = at.Add(upstreamInterval)
	blocked := c.blockedUntil
	c.mu.Unlock()

	if now.Before(blocked) {
		return errThrottled
	}
	if d := time.Until(at); d > 0 {
		t := time.NewTimer(d)
		defer t.Stop()
		select {
		case <-t.C:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

func (c *Client) throttle() {
	c.mu.Lock()
	c.blockedUntil = time.Now().Add(throttleCooldown)
	c.mu.Unlock()
}

type summaryResponse struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Extract     string `json:"extract"`
	Thumbnail   *struct {
		Source string `json:"source"`
	} `json:"thumbnail"`
	ContentURLs struct {
		Desktop struct {
			Page string `json:"page"`
		} `json:"desktop"`
	} `json:"content_urls"`
	Titles struct {
		Display string `json:"display"`
	} `json:"titles"`
}

func (c *Client) Search(ctx context.Context, query, language string) (*Article, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, nil
	}

	lang := primaryLanguage(query, language)
	return c.lookup(ctx, query, lang)
}

func (c *Client) lookup(ctx context.Context, query, language string) (*Article, error) {
	key := language + "\x00" + query
	if art, ok := c.fromCache(key); ok {
		return art, nil
	}
	if err := c.waitSlot(ctx); err != nil {
		return nil, err
	}

	lctx, cancel := context.WithTimeout(ctx, perLanguageTimeout)
	defer cancel()

	art, err := c.fetchSummary(lctx, query, language)
	if err == nil && lctx.Err() == nil {
		c.toCache(key, art)
	}
	return art, err
}

func (c *Client) fromCache(key string) (*Article, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.cache[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(entry.expires) {
		delete(c.cache, key)
		return nil, false
	}
	return entry.article, true
}

func (c *Client) toCache(key string, article *Article) {
	ttl := cacheTTL
	if article == nil {
		ttl = negativeCacheTTL
	}
	expires := time.Now().Add(ttl)

	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.cache) >= maxCacheEntries {
		now := time.Now()
		for k, v := range c.cache {
			if now.After(v.expires) {
				delete(c.cache, k)
			}
		}
		if len(c.cache) >= maxCacheEntries {
			c.cache = make(map[string]cacheEntry, maxCacheEntries)
		}
	}
	c.cache[key] = cacheEntry{article: article, expires: expires}
}

func (c *Client) fetchSummary(ctx context.Context, query, language string) (*Article, error) {
	rawURL := fmt.Sprintf(summaryURL, language, url.PathEscape(wikiTitle(query)))

	for redirect := 0; ; redirect++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Accept-Language", language)
		req.Header.Set("User-Agent", userAgent)

		res, err := c.http.Do(req)
		if err != nil {
			return nil, err
		}

		switch {
		case res.StatusCode == http.StatusOK:
			defer res.Body.Close()
			var payload summaryResponse
			if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
				return nil, err
			}
			return articleFromSummary(payload, language), nil

		case res.StatusCode == http.StatusTooManyRequests:
			res.Body.Close()
			c.throttle()
			return nil, errThrottled

		case res.StatusCode >= 300 && res.StatusCode < 400 && redirect < maxRedirects:
			loc, locErr := res.Location()
			_, _ = io.Copy(io.Discard, res.Body)
			res.Body.Close()
			if locErr != nil {
				return nil, nil
			}
			rawURL = loc.String()
			continue

		case res.StatusCode == http.StatusNotFound, res.StatusCode == http.StatusBadRequest:
			res.Body.Close()
			return nil, nil

		default:
			status := res.Status
			res.Body.Close()
			return nil, fmt.Errorf("wikipedia: unexpected status %s", status)
		}
	}
}

func articleFromSummary(p summaryResponse, language string) *Article {
	if p.Type != "standard" {
		return nil
	}

	title := strings.TrimSpace(helpers.CleanText(p.Titles.Display))
	if title == "" {
		title = p.Title
	}
	if title == "" {
		return nil
	}

	pageURL := p.ContentURLs.Desktop.Page
	if pageURL == "" {
		pageURL = fmt.Sprintf("https://%s.wikipedia.org/wiki/%s", language,
			url.PathEscape(strings.ReplaceAll(title, " ", "_")))
	}

	article := &Article{
		Title:       title,
		Description: strings.TrimSpace(p.Description),
		Extract:     strings.TrimSpace(p.Extract),
		URL:         pageURL,
	}
	if p.Thumbnail != nil {
		article.Thumbnail = p.Thumbnail.Source
	}
	return article
}

func wikiTitle(query string) string {
	if !isLower(query) {
		return query
	}
	return titleCase(query)
}

func isLower(s string) bool {
	hasCased := false
	for _, r := range s {
		if unicode.IsUpper(r) || unicode.IsTitle(r) {
			return false
		}
		if unicode.IsLower(r) {
			hasCased = true
		}
	}
	return hasCased
}

func titleCase(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	inWord := false
	for _, r := range s {
		if unicode.IsLetter(r) {
			if inWord {
				b.WriteRune(unicode.ToLower(r))
			} else {
				b.WriteRune(unicode.ToUpper(r))
			}
			inWord = true
			continue
		}
		b.WriteRune(r)
		inWord = false
	}
	return b.String()
}

var languagePattern = regexp.MustCompile(`^[a-z]{2,3}$`)

func primaryLanguage(query, requested string) string {
	requested = strings.ToLower(strings.TrimSpace(requested))
	if requested != "" && requested != "auto" && languagePattern.MatchString(requested) {
		return requested
	}

	for _, r := range strings.ToLower(query) {
		switch {
		case r >= '\u0400' && r <= '\u04FF': // cyrillic
			if strings.ContainsAny(query, "іїєґ") {
				return "uk"
			}
			return "ru"
		case r >= '\u4E00' && r <= '\u9FFF', r >= '\u3400' && r <= '\u4DBF': // CJK Unified Ideographs
			return "zh"
		case r >= '\u3040' && r <= '\u309F', r >= '\u30A0' && r <= '\u30FF': // Hiragana + Katakana
			return "ja"
		case r >= '\uAC00' && r <= '\uD7AF': // hangul syllables
			return "ko"
		case r >= '\u0600' && r <= '\u06FF': // arabic
			return "ar"
		case r >= '\u0900' && r <= '\u097F': // devanagari
			return "hi"
		case r >= '\u0E00' && r <= '\u0E5F': // thai
			return "th"
		}
	}

	return "en"
}
