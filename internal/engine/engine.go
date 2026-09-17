package engine

import (
	"context"
	"net/http"
)

const (
	// BrowserUserAgent is the shared User-Agent string for all engine requests
	BrowserUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/150.0.0.0 Safari/537.36"

	// MaxBodyReadLimit is the default maximum bytes to read from an HTTP response body
	MaxBodyReadLimit = 8 << 20 // 8 MB
)

type Sitelink struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type Result struct {
	Title     string     `json:"title"`
	URL       string     `json:"url"`
	Snippet   string     `json:"snippet"`
	Engine    string     `json:"engine"`
	Engines   []string   `json:"engines,omitempty"`
	Sitelinks []Sitelink `json:"sitelinks,omitempty"`
}

type Infobox struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type ImageResult struct {
	Title     string   `json:"title"`
	URL       string   `json:"url"`
	SourceURL string   `json:"source_url,omitempty"`
	Thumbnail string   `json:"thumbnail"`
	Width     int      `json:"width"`
	Height    int      `json:"height"`
	Engine    string   `json:"engine"`
	Engines   []string `json:"engines,omitempty"`
}

// abstracts a search backend for the aggregator
type Engine interface {
	Name() string
	Search(ctx context.Context, query string, limit, offset int) ([]Result, error)
}

type ImageSearcher interface {
	ImageSearch(ctx context.Context, query string, limit, offset int) ([]ImageResult, error)
}

// minimal *http.Client subset so engines can be tested with fake transports
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}
