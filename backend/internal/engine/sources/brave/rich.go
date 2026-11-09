package brave

import (
	"encoding/json"
	"strings"

	"clasp/internal/engine"
	"clasp/internal/engine/helpers"
)

const richMarker = `{type:"rich",results:[`

func parseRich(doc string) []engine.RichAnswer {
	i := strings.Index(doc, richMarker)
	if i < 0 {
		return nil
	}
	i += len(richMarker) - 1
	end, ok := helpers.MatchObject(doc, i)
	if !ok {
		return nil
	}
	v, err := helpers.ParseJSLiteral(doc[i:end])
	if err != nil {
		return nil
	}
	arr, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]engine.RichAnswer, 0, len(arr))
	for _, item := range arr {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if ra := richAnswerFromMap(m); ra != nil {
			out = append(out, *ra)
		}
	}
	return out
}

type richSubtype struct {
	name    string
	aliases []string
}

var richSubtypes = []richSubtype{
	{name: "cryptocurrency"},
	{name: "stock", aliases: []string{"stocks"}},
}

func (s richSubtype) matches(subtype string) bool {
	if subtype == "" || subtype == s.name {
		return true
	}
	for _, a := range s.aliases {
		if subtype == a {
			return true
		}
	}
	return false
}

func richAnswerFromMap(m map[string]any) *engine.RichAnswer {
	f := fields(m)
	subtype := strings.ToLower(f.str("subtype"))
	for _, st := range richSubtypes {
		if !st.matches(subtype) {
			continue
		}
		payload := f.obj(st.name)
		if len(payload) == 0 {
			continue
		}
		if a := richAnswerFor(st.name, payload); a != nil {
			return a
		}
		if subtype != "" {
			break
		}
	}
	return nil
}

func richAnswerFor(name string, payload fields) *engine.RichAnswer {
	switch name {
	case "cryptocurrency":
		return richCryptocurrency(payload)
	case "stock":
		return richStock(payload)
	}
	return nil
}

func richCryptocurrency(p fields) *engine.RichAnswer {
	q := p.obj("quote")
	if q == nil {
		return nil
	}
	out := &engine.RichCryptocurrency{IntentType: p.str("intent_type")}
	out.Quote = &engine.RichCryptoQuote{
		ID:                       q.str("id"),
		Symbol:                   q.str("symbol"),
		Name:                     q.str("name"),
		Image:                    q.str("image"),
		CurrentPrice:             q.num("current_price"),
		MarketCap:                q.num("market_cap"),
		High24h:                  q.num("high_24h"),
		Low24h:                   q.num("low_24h"),
		PriceChange24h:           q.num("price_change_24h"),
		PriceChangePercentage24h: q.num("price_change_percentage_24h"),
	}
	if ts := p.obj("timeseries"); ts != nil {
		vs := ts.str("vs_currency")
		if vs == "" {
			vs = p.str("vs_currency")
		}
		out.VsCurrency = vs
		if series := richPriceSeries(ts); len(series) > 1 {
			out.Timeseries = &engine.RichCryptoSeries{
				TimeRange: ts.str("time_range"),
				TsPrice:   series,
			}
		}
	}
	return &engine.RichAnswer{Type: "cryptocurrency", Cryptocurrency: out}
}

func richStock(p fields) *engine.RichAnswer {
	q := p.obj("quote")
	if q == nil {
		return nil
	}
	out := &engine.RichStock{TimeRange: p.str("time_range")}
	if ai := p.obj("asset_info"); ai != nil {
		out.AssetInfo = &engine.RichStockAsset{
			Symbol:   ai.str("symbol"),
			Name:     ai.str("name"),
			Category: ai.str("category"),
			Currency: ai.str("currency"),
			Country:  ai.str("country"),
			Exchange: ai.str("exchange"),
			MicCode:  ai.str("mic_code"),
		}
	}
	if xi := p.obj("exchange_info"); xi != nil {
		out.ExchangeInfo = &engine.RichStockExchange{
			Timezone:  xi.str("timezone"),
			OpenTime:  xi.str("open_time"),
			CloseTime: xi.str("close_time"),
		}
	}
	out.Quote = &engine.RichStockQuote{
		Symbol:          q.str("symbol"),
		CompanyName:     q.str("company_name"),
		PrimaryExchange: q.str("primary_exchange"),
		Currency:        q.str("currency"),
		LatestPrice:     q.num("latest_price"),
		Open:            q.num("open"),
		High:            q.num("high"),
		Low:             q.num("low"),
		Close:           q.num("close"),
		Change:          q.num("change"),
		ChangePercent:   q.num("change_percent"),
		Volume:          q.num("volume"),
		MarketCap:       q.num("market_cap"),
		PeRatio:         q.num("pe_ratio"),
		Week52High:      q.num("week_52_high"),
		Week52Low:       q.num("week_52_low"),
		LatestUpdate:    q.i64("latest_update"),
	}
	if ts := p.obj("timeseries"); ts != nil {
		if series := richStockSeries(ts); len(series) > 1 {
			out.Timeseries = &engine.RichStockSeries{
				TimeRange: ts.str("time_range"),
				Close:     series,
			}
		}
	}
	return &engine.RichAnswer{Type: "stock", Stock: out}
}

func richStockSeries(ts fields) [][2]float64 {
	raw := ts.list("timeseries")
	out := make([][2]float64, 0, len(raw))
	for _, r := range raw {
		m, ok := r.(map[string]any)
		if !ok {
			continue
		}
		f := fields(m)
		t, tok := richFloat(f["timestamp"])
		v, vok := richFloat(f["close"])
		if !tok || !vok {
			continue
		}
		out = append(out, [2]float64{t, v})
	}
	return downsampleSeries(out, maxSeriesPoints)
}

func downsampleSeries(out [][2]float64, maxPoints int) [][2]float64 {
	if len(out) <= maxPoints {
		return out
	}
	stride := (len(out) + maxPoints - 1) / maxPoints
	sampled := make([][2]float64, 0, maxPoints)
	for i := 0; i < len(out); i += stride {
		sampled = append(sampled, out[i])
	}
	return sampled
}

const maxSeriesPoints = 96

func richPriceSeries(ts fields) [][2]float64 {
	raw := ts.list("ts_price")
	out := make([][2]float64, 0, len(raw))
	for _, r := range raw {
		pair, ok := r.([]any)
		if !ok || len(pair) < 2 {
			continue
		}
		t, tok := richFloat(pair[0])
		v, vok := richFloat(pair[1])
		if !tok || !vok {
			continue
		}
		out = append(out, [2]float64{t, v})
	}
	return downsampleSeries(out, maxSeriesPoints)
}

func richFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	default:
		if i, ok := intVal(n); ok {
			return float64(i), true
		}
	}
	return 0, false
}

func intVal(v any) (int64, bool) {
	switch n := v.(type) {
	case int:
		return int64(n), true
	case int64:
		return n, true
	case float64:
		return int64(n), true
	}
	return 0, false
}
