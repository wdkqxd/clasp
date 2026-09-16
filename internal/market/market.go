package market

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type Point struct {
	Time  int64   `json:"time"`
	Value float64 `json:"value"`
}

type Quote struct {
	Query         string  `json:"query"`
	Symbol        string  `json:"symbol"`
	Name          string  `json:"name"`
	Currency      string  `json:"currency"`
	Current       float64 `json:"current"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"change_percent"`
	Period        string  `json:"period"`
	Points        []Point `json:"points"`
}

type Client struct {
	http *http.Client
}

var pairPattern = regexp.MustCompile(`(?i)^([a-z]{3})\s*(?:/|to|in|в)\s*([a-z]{3})$`)
var tickerPattern = regexp.MustCompile(`^[A-Z][A-Z0-9.-]{0,9}$`)

var cryptoIDs = map[string]string{
	"BTC": "bitcoin", "ETH": "ethereum", "SOL": "solana", "XRP": "ripple",
	"ADA": "cardano", "DOGE": "dogecoin", "DOT": "polkadot", "AVAX": "avalanche-2",
	"LINK": "chainlink", "MATIC": "matic-network", "LTC": "litecoin", "BCH": "bitcoin-cash",
	"TON": "the-open-network", "TRX": "tron", "SHIB": "shiba-inu", "XLM": "stellar",
}

var currencyCodes = func() map[string]bool {
	codes := strings.Fields("AED AFN ALL AMD ANG AOA ARS AUD AWG AZN BAM BBD BDT BGN BHD BIF BMD BND BOB BOV BRL BSD BTN BWP BYN BZD CAD CDF CHE CHF CHW CLF CLP CNY COP COU CRC CUC CUP CVE CZK DJF DKK DOP DZD EGP ERN ETB EUR FJD FKP GBP GEL GHS GIP GMD GNF GTQ GYD HKD HNL HTG HUF IDR ILS INR IQD IRR ISK JMD JOD JPY KES KGS KHR KMF KPW KRW KWD KYD KZT LAK LBP LKR LRD LSL LYD MAD MDL MGA MKD MMK MNT MOP MRU MUR MVR MWK MXN MXV MYR MZN NAD NGN NIO NOK NPR NZD OMR PAB PEN PGK PHP PKR PLN PYG QAR RON RSD RUB RWF SAR SBD SCR SDG SEK SGD SHP SLE SLL SOS SRD SSP STN SVC SYP SZL THB TJS TMT TND TOP TRY TTD TWD TZS UAH UGX USD USN UYI UYU UYW UZS VED VES VND VUV WST XAF XAG XAU XBA XBB XBC XBD XCD XDR XOF XPD XPF XPT XSU XTS XUA XXX YER ZAR ZMW ZWL")
	result := make(map[string]bool, len(codes))
	for _, code := range codes {
		result[code] = true
	}
	return result
}()

func NewClient() *Client { return &Client{http: &http.Client{Timeout: 12 * time.Second}} }

func (c *Client) Search(ctx context.Context, rawQuery, period string) (*Quote, error) {
	q := strings.TrimSpace(rawQuery)
	if q == "" {
		return nil, nil
	}
	if period != "1d" && period != "7d" && period != "1m" && period != "1y" {
		period = "7d"
	}
	if m := pairPattern.FindStringSubmatch(q); len(m) == 3 && currencyCodes[strings.ToUpper(m[1])] && currencyCodes[strings.ToUpper(m[2])] {
		return c.currency(ctx, q, strings.ToUpper(m[1]), strings.ToUpper(m[2]), period)
	}

	symbol := strings.ToUpper(q)
	if id, ok := cryptoIDs[symbol]; ok {
		return c.coingecko(ctx, q, symbol, id, period)
	}

	if currencyCodes[symbol] {
		return nil, nil
	}
	if !tickerPattern.MatchString(symbol) || strings.Contains(q, " ") {
		return nil, nil
	}
	return c.yahoo(ctx, q, symbol, period)
}

func (c *Client) requestJSON(ctx context.Context, endpoint string, target any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "metasearch/1.0")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("market provider returned %s", resp.Status)
	}
	return json.NewDecoder(io.LimitReader(resp.Body, 8<<20)).Decode(target)
}

type cgResponse struct {
	Prices [][]float64 `json:"prices"`
}

func (c *Client) coingecko(ctx context.Context, query, symbol, id, period string) (*Quote, error) {
	days := map[string]string{"1d": "1", "7d": "7", "1m": "30", "1y": "365"}[period]
	endpoint := "https://api.coingecko.com/api/v3/coins/" + url.PathEscape(id) + "/market_chart?vs_currency=usd&days=" + days
	var data cgResponse
	if err := c.requestJSON(ctx, endpoint, &data); err != nil || len(data.Prices) < 2 {
		return nil, err
	}
	points := make([]Point, 0, len(data.Prices))
	for _, item := range data.Prices {
		if len(item) == 2 && finite(item[1]) {
			points = append(points, Point{Time: int64(item[0]), Value: item[1]})
		}
	}
	return quoteFromPoints(query, symbol, symbol, "USD", period, points), nil
}

type yahooResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Currency      string  `json:"currency"`
				ShortName     string  `json:"shortName"`
				QuoteType     string  `json:"quoteType"`
				Instrument    string  `json:"instrumentType"`
				Exchange      string  `json:"exchangeName"`
				RegularMarket float64 `json:"regularMarketPrice"`
			} `json:"meta"`
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Close []*float64 `json:"close"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
	} `json:"chart"`
}

func (c *Client) yahoo(ctx context.Context, query, symbol, period string) (*Quote, error) {
	ranges := map[string][2]string{"1d": {"1d", "15m"}, "7d": {"7d", "1h"}, "1m": {"1mo", "1d"}, "1y": {"1y", "1d"}}
	r := ranges[period]
	endpoint := "https://query1.finance.yahoo.com/v8/finance/chart/" + url.PathEscape(symbol) + "?range=" + r[0] + "&interval=" + r[1]
	var data yahooResponse
	if err := c.requestJSON(ctx, endpoint, &data); err != nil || len(data.Chart.Result) == 0 {
		return nil, err
	}
	result := data.Chart.Result[0]
	kind := strings.ToUpper(result.Meta.QuoteType + " " + result.Meta.Instrument + " " + result.Meta.Exchange + " " + result.Meta.ShortName)
	if strings.Contains(kind, "CRYPTO") || strings.Contains(kind, "TOKEN") {
		return nil, nil
	}
	if len(result.Indicators.Quote) == 0 {
		return nil, nil
	}
	closes := result.Indicators.Quote[0].Close
	points := make([]Point, 0, len(closes))
	for i, value := range closes {
		if value != nil && i < len(result.Timestamp) && finite(*value) {
			points = append(points, Point{Time: result.Timestamp[i] * 1000, Value: *value})
		}
	}
	if len(points) < 2 {
		return nil, nil
	}
	name := result.Meta.ShortName
	if name == "" {
		name = symbol
	}
	currency := strings.ToUpper(result.Meta.Currency)
	if currency == "" {
		currency = "USD"
	}
	return quoteFromPoints(query, symbol, name, currency, period, points), nil
}

type frankfurterResponse struct {
	Amount float64            `json:"amount"`
	Base   string             `json:"base"`
	Date   string             `json:"date"`
	Rates  map[string]float64 `json:"rates"`
}

func (c *Client) currency(ctx context.Context, query, from, to, period string) (*Quote, error) {
	days := map[string]int{"1d": 1, "7d": 7, "1m": 30, "1y": 365}[period]
	end := time.Now().UTC()
	start := end.AddDate(0, 0, -days)
	endpoint := "https://api.frankfurter.dev/v2/rates?base=" + from + "&symbols=" + to + "&from=" + start.Format("2006-01-02") + "&to=" + end.Format("2006-01-02")
	var data []frankfurterResponse
	if err := c.requestJSON(ctx, endpoint, &data); err != nil {
		return nil, err
	}
	points := make([]Point, 0, len(data)+1)
	for _, item := range data {
		if value, ok := item.Rates[to]; ok && finite(value) {
			t, err := time.Parse("2006-01-02", item.Date)
			if err == nil {
				points = append(points, Point{Time: t.UnixMilli(), Value: value})
			}
		}
	}
	if len(points) < 2 {
		return nil, nil
	}
	return quoteFromPoints(query, from+"/"+to, from+" → "+to, to, period, points), nil
}

func quoteFromPoints(query, symbol, name, currency, period string, points []Point) *Quote {
	first := points[0].Value
	last := points[len(points)-1].Value
	change := last - first
	percent := 0.0
	if first != 0 {
		percent = change / math.Abs(first) * 100
	}
	return &Quote{Query: query, Symbol: symbol, Name: name, Currency: currency, Current: last, Change: change, ChangePercent: percent, Period: period, Points: points}
}

func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }
