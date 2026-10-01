package suggest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"metasearch/internal/engine"
)

const (
	EngineBrave     = "brave"
	EngineStartpage = "startpage"
	EngineDDG       = "ddg"
	EngineKagi      = "kagi"

	DefaultEngine = EngineKagi

	requestTimeout = 4 * time.Second
	cacheTTL       = 5 * time.Minute
	maxCacheItems  = 512
)

const MaxSuggestions = 7

type Response struct {
	Query       string   `json:"query"`
	Engine      string   `json:"engine"`
	Suggestions []string `json:"suggestions"`
}

func Engines() []string {
	return []string{EngineKagi, EngineDDG, EngineBrave, EngineStartpage}
}

func Normalize(name string) string {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "":
		return DefaultEngine
	case EngineBrave:
		return EngineBrave
	case EngineStartpage:
		return EngineStartpage
	case EngineDDG, "duckduckgo":
		return EngineDDG
	case EngineKagi:
		return EngineKagi
	default:
		return ""
	}
}

type cacheEntry struct {
	suggestions []string
	expires     time.Time
}

type Client struct {
	http engine.HTTPClient

	mu    sync.Mutex
	cache map[string]cacheEntry
}

func NewClient(httpClient engine.HTTPClient) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: requestTimeout}
	}
	return &Client{http: httpClient, cache: make(map[string]cacheEntry)}
}

func (c *Client) Suggest(ctx context.Context, query, engineName string) (Response, error) {
	name := Normalize(engineName)
	if name == "" {
		return Response{}, fmt.Errorf("unsupported autocomplete engine: %q", engineName)
	}
	query = strings.TrimSpace(query)
	if query == "" {
		return Response{Engine: name, Suggestions: []string{}}, nil
	}

	if cached, ok := c.cached(name, query); ok {
		return Response{Query: query, Engine: name, Suggestions: cached}, nil
	}

	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	var (
		suggestions []string
		err         error
	)
	switch name {
	case EngineBrave:
		suggestions, err = c.fetchArray(ctx, "https://search.brave.com/api/suggest", query, url.Values{"country": {"ALL"}})
	case EngineStartpage:
		suggestions, err = c.fetchArray(ctx, "https://www.startpage.com/osuggestions", query, nil)
	case EngineDDG:
		suggestions, err = c.fetchPhrases(ctx, "https://duckduckgo.com/ac/", query, url.Values{"kl": {"wt-wt"}})
	case EngineKagi:
		suggestions, err = c.fetchTitles(ctx, "https://kagi.com/autosuggest", query, nil)
	}
	if err != nil {
		return Response{}, err
	}

	suggestions = normalizeSuggestions(query, suggestions)
	c.store(name, query, suggestions)
	return Response{Query: query, Engine: name, Suggestions: suggestions}, nil
}

func (c *Client) fetchArray(ctx context.Context, endpoint, query string, extra url.Values) ([]string, error) {
	body, err := c.get(ctx, endpoint, query, extra)
	if err != nil {
		return nil, err
	}
	var payload []json.RawMessage
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	if len(payload) < 2 {
		return nil, nil
	}
	var list []string
	if err := json.Unmarshal(payload[1], &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (c *Client) fetchPhrases(ctx context.Context, endpoint, query string, extra url.Values) ([]string, error) {
	body, err := c.get(ctx, endpoint, query, extra)
	if err != nil {
		return nil, err
	}
	var payload []struct {
		Phrase string `json:"phrase"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	list := make([]string, 0, len(payload))
	for _, item := range payload {
		if item.Phrase != "" {
			list = append(list, item.Phrase)
		}
	}
	return list, nil
}

func (c *Client) fetchTitles(ctx context.Context, endpoint, query string, extra url.Values) ([]string, error) {
	body, err := c.get(ctx, endpoint, query, extra)
	if err != nil {
		return nil, err
	}
	var payload []struct {
		Title string `json:"t"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	list := make([]string, 0, len(payload))
	for _, item := range payload {
		if item.Title != "" {
			list = append(list, item.Title)
		}
	}
	return list, nil
}

func (c *Client) get(ctx context.Context, endpoint, query string, extra url.Values) ([]byte, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("q", query)
	for key, values := range extra {
		for _, value := range values {
			q.Set(key, value)
		}
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", engine.BrowserUserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := c.http.Do(req)
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
	return io.ReadAll(io.LimitReader(resp.Body, engine.MaxBodyReadLimit))
}

func (c *Client) cached(name, query string) ([]string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.cache[name+"\x00"+strings.ToLower(query)]
	if !ok || time.Now().After(entry.expires) {
		return nil, false
	}
	return entry.suggestions, true
}

func (c *Client) store(name, query string, suggestions []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.cache) >= maxCacheItems {
		c.cache = make(map[string]cacheEntry)
	}
	c.cache[name+"\x00"+strings.ToLower(query)] = cacheEntry{
		suggestions: suggestions,
		expires:     time.Now().Add(cacheTTL),
	}
}

func normalizeSuggestions(query string, suggestions []string) []string {
	seen := make(map[string]struct{}, len(suggestions))
	out := make([]string, 0, len(suggestions))
	for _, s := range suggestions {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		key := strings.ToLower(s)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, s)
		if len(out) >= MaxSuggestions {
			break
		}
	}
	return out
}
