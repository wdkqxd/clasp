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

type Discussion struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet,omitempty"`
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

type RichAnswer struct {
	Type           string              `json:"type"`
	Cryptocurrency *RichCryptocurrency `json:"cryptocurrency,omitempty"`
	Stock          *RichStock          `json:"stock,omitempty"`
}

type RichCryptocurrency struct {
	IntentType string            `json:"intent_type"`
	VsCurrency string            `json:"vs_currency"`
	Quote      *RichCryptoQuote  `json:"quote,omitempty"`
	Timeseries *RichCryptoSeries `json:"timeseries,omitempty"`
}

type RichCryptoQuote struct {
	ID                       string  `json:"id"`
	Symbol                   string  `json:"symbol"`
	Name                     string  `json:"name"`
	Image                    string  `json:"image,omitempty"`
	CurrentPrice             float64 `json:"current_price"`
	MarketCap                float64 `json:"market_cap"`
	High24h                  float64 `json:"high_24h"`
	Low24h                   float64 `json:"low_24h"`
	PriceChange24h           float64 `json:"price_change_24h"`
	PriceChangePercentage24h float64 `json:"price_change_percentage_24h"`
}

type RichCryptoSeries struct {
	TimeRange string       `json:"time_range"`
	TsPrice   [][2]float64 `json:"ts_price,omitempty"`
}

type RichStock struct {
	AssetInfo    *RichStockAsset    `json:"asset_info,omitempty"`
	ExchangeInfo *RichStockExchange `json:"exchange_info,omitempty"`
	TimeRange    string             `json:"time_range"`
	Quote        *RichStockQuote    `json:"quote,omitempty"`
	Timeseries   *RichStockSeries   `json:"timeseries,omitempty"`
}

type RichStockAsset struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Category string `json:"category,omitempty"`
	Currency string `json:"currency,omitempty"`
	Country  string `json:"country,omitempty"`
	Exchange string `json:"exchange,omitempty"`
	MicCode  string `json:"mic_code,omitempty"`
}

type RichStockExchange struct {
	Timezone  string `json:"timezone,omitempty"`
	OpenTime  string `json:"open_time,omitempty"`
	CloseTime string `json:"close_time,omitempty"`
}

type RichStockQuote struct {
	Symbol          string  `json:"symbol"`
	CompanyName     string  `json:"company_name,omitempty"`
	PrimaryExchange string  `json:"primary_exchange,omitempty"`
	Currency        string  `json:"currency,omitempty"`
	LatestPrice     float64 `json:"latest_price"`
	Open            float64 `json:"open"`
	High            float64 `json:"high"`
	Low             float64 `json:"low"`
	Close           float64 `json:"close"`
	Change          float64 `json:"change"`
	ChangePercent   float64 `json:"change_percent"`
	Volume          float64 `json:"volume"`
	MarketCap       float64 `json:"market_cap,omitempty"`
	PeRatio         float64 `json:"pe_ratio,omitempty"`
	Week52High      float64 `json:"week_52_high,omitempty"`
	Week52Low       float64 `json:"week_52_low,omitempty"`
	LatestUpdate    int64   `json:"latest_update,omitempty"`
}

type RichStockSeries struct {
	TimeRange string       `json:"time_range"`
	Close     [][2]float64 `json:"ts_close,omitempty"`
}

type SearchResult struct {
	Results     []Result
	Infobox     *Infobox
	Discussions []Discussion
	Rich        []RichAnswer
}

// abstracts a search backend for the aggregator
type Engine interface {
	Name() string
	Search(ctx context.Context, query string, limit, offset int) (SearchResult, error)
}

type ImageSearcher interface {
	ImageSearch(ctx context.Context, query string, limit, offset int) ([]ImageResult, error)
}

// minimal *http.Client subset so engines can be tested with fake transports
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}
