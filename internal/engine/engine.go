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
	Calculator     *RichCalculator     `json:"calculator,omitempty"`
	Colorpicker    *RichColorpicker    `json:"colorpicker,omitempty"`
	Weather        *RichWeather        `json:"weather,omitempty"`
	Timer          *RichTimer          `json:"timer,omitempty"`
	UnitConversion *RichUnitConversion `json:"unitConversion,omitempty"`
	Timezones      *RichTimezones      `json:"timezones,omitempty"`
	Cryptocurrency *RichCryptocurrency `json:"cryptocurrency,omitempty"`
	Unixtimestamp  *RichUnixtimestamp  `json:"unixtimestamp,omitempty"`
	Currency       *RichCurrency       `json:"currency,omitempty"`
	News           *RichNews           `json:"news,omitempty"`
	TopAnswer      *RichTopAnswer      `json:"topAnswer,omitempty"`
	Stopwatch      bool                `json:"stopwatch,omitempty"`
}

type RichTopAnswer struct {
	Question string            `json:"question"`
	Blocks   []RichAnswerBlock `json:"blocks"`
	Author   string            `json:"author,omitempty"`
	Upvotes  int               `json:"upvotes,omitempty"`
	Title    string            `json:"title,omitempty"`
	URL      string            `json:"url"`
	MetaURL  *RichMetaURL      `json:"meta_url,omitempty"`
}

type RichAnswerBlock struct {
	Kind  string   `json:"kind"`
	Text  string   `json:"text,omitempty"`
	Items []string `json:"items,omitempty"`
}

type RichCalculator struct {
	Expression string  `json:"expression"`
	Answer     float64 `json:"answer"`
}

type RichColorpicker struct {
	Color  string `json:"color"`
	Format string `json:"format"`
}

type RichWeather struct {
	Location       *RichWeatherLocation `json:"location,omitempty"`
	CurrentWeather *RichWeatherCurrent  `json:"current_weather,omitempty"`
	Daily          []RichWeatherDay     `json:"daily,omitempty"`
	Alerts         []RichWeatherAlert   `json:"alerts,omitempty"`
}

type RichWeatherLocation struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	State   string `json:"state"`
}

type RichWeatherCond struct {
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type RichWeatherCurrent struct {
	Temp      float64          `json:"temp"`
	FeelsLike float64          `json:"feels_like"`
	Humidity  float64          `json:"humidity"`
	WindSpeed float64          `json:"wind_speed"`
	UVI       float64          `json:"uvi"`
	Pressure  float64          `json:"pressure"`
	Weather   *RichWeatherCond `json:"weather,omitempty"`
}

type RichWeatherTemp struct {
	Max float64 `json:"max"`
	Min float64 `json:"min"`
	Day float64 `json:"day"`
}

type RichWeatherDay struct {
	Ts          int64            `json:"ts"`
	DateI18n    string           `json:"date_i18n"`
	Temperature *RichWeatherTemp `json:"temperature,omitempty"`
	Weather     *RichWeatherCond `json:"weather,omitempty"`
}

type RichWeatherAlert struct {
	Event             string `json:"event"`
	Description       string `json:"description,omitempty"`
	StartRelativeI18n string `json:"start_relative_i18n,omitempty"`
}

type RichTimer struct {
	Duration    int64 `json:"duration"`
	StartOnLoad bool  `json:"start_on_load"`
}

type RichUnitConversion struct {
	Amount         float64 `json:"amount"`
	FromUnit       string  `json:"from_unit"`
	ToUnit         string  `json:"to_unit"`
	Dimensionality string  `json:"dimensionality"`
}

type RichTimezones struct {
	Type      string         `json:"type"`
	Timezones []RichTimezone `json:"timezones"`
}

type RichTimezone struct {
	Abbreviation  string           `json:"abbreviation"`
	GenericName   string           `json:"generic_name,omitempty"`
	UTCOffset     string           `json:"utc_offset,omitempty"`
	ConvertedTime *RichTimezoneCnv `json:"converted_time,omitempty"`
}

type RichTimezoneCnv struct {
	Strftime string `json:"strftime,omitempty"`
	Location string `json:"location,omitempty"`
	UTCDiff  string `json:"utc_diff,omitempty"`
	City     string `json:"city,omitempty"`
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

type RichUnixtimestamp struct {
	Conversion any `json:"conversion,omitempty"`
}

type RichCurrency struct {
	Amount       float64         `json:"amount"`
	FromCode     string          `json:"from_currency_code"`
	ToCode       string          `json:"to_currency_code"`
	FromName     string          `json:"from_currency_name"`
	ToName       string          `json:"to_currency_name"`
	Result       float64         `json:"result"`
	ExchangeRate float64         `json:"exchange_rate"`
	Date         string          `json:"date"`
	Timeseries   *RichCurrSeries `json:"timeseries,omitempty"`
}

type RichCurrSeries struct {
	TimeRange  string   `json:"time_range"`
	TsExchange [][2]any `json:"ts_exchange_rate,omitempty"`
}

type RichNews struct {
	IsTopical bool              `json:"is_topical"`
	Topics    []string          `json:"topics,omitempty"`
	Articles  []RichNewsArticle `json:"articles"`
}

type RichNewsArticle struct {
	Title         string       `json:"title"`
	URL           string       `json:"url"`
	Description   string       `json:"description,omitempty"`
	Img           string       `json:"img,omitempty"`
	PublishTime   int64        `json:"publish_time,omitempty"`
	PublisherName string       `json:"publisher_name,omitempty"`
	PageAge       string       `json:"page_age,omitempty"`
	MetaURL       *RichMetaURL `json:"meta_url,omitempty"`
}

type RichMetaURL struct {
	Hostname string `json:"hostname,omitempty"`
	Favicon  string `json:"favicon,omitempty"`
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
