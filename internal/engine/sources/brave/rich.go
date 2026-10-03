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
	i += len(richMarker) - 1 // point at the opening '['
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

var richSubtypes = []string{"cryptocurrency"}

func richAnswerFromMap(m map[string]any) *engine.RichAnswer {
	subtype, _ := m["subtype"].(string)
	subtype = strings.ToLower(subtype)
	for _, name := range richSubtypes {
		if subtype != "" && subtype != name {
			continue
		}
		payload, ok := m[name].(map[string]any)
		if !ok || len(payload) == 0 {
			continue
		}
		if a := richAnswerFor(name, payload); a != nil {
			return a
		}
		// subtype matched but payload not usable -> stop scanning
		if subtype == name {
			break
		}
	}
	return nil
}

func richAnswerFor(name string, payload map[string]any) *engine.RichAnswer {
	switch name {
	case "cryptocurrency":
		return richCryptocurrency(payload)
	}
	return nil
}

func richCryptocurrency(p map[string]any) *engine.RichAnswer {
	q, ok := p["quote"].(map[string]any)
	if !ok {
		return nil
	}
	out := &engine.RichCryptocurrency{IntentType: richStr(p["intent_type"])}
	out.Quote = &engine.RichCryptoQuote{
		ID:                       richStr(q["id"]),
		Symbol:                   richStr(q["symbol"]),
		Name:                     richStr(q["name"]),
		Image:                    richStr(q["image"]),
		CurrentPrice:             richFloatOrDefault(q["current_price"]),
		MarketCap:                richFloatOrDefault(q["market_cap"]),
		High24h:                  richFloatOrDefault(q["high_24h"]),
		Low24h:                   richFloatOrDefault(q["low_24h"]),
		PriceChange24h:           richFloatOrDefault(q["price_change_24h"]),
		PriceChangePercentage24h: richFloatOrDefault(q["price_change_percentage_24h"]),
	}
	if ts, ok := p["timeseries"].(map[string]any); ok {
		vs := richStr(ts["vs_currency"])
		if vs == "" {
			vs = richStr(p["vs_currency"])
		}
		out.VsCurrency = vs
		if series := richPriceSeries(ts); len(series) > 1 {
			out.Timeseries = &engine.RichCryptoSeries{
				TimeRange: richStr(ts["time_range"]),
				TsPrice:   series,
			}
		}
	}
	return &engine.RichAnswer{Type: "cryptocurrency", Cryptocurrency: out}
}

func richPriceSeries(ts map[string]any) [][2]float64 {
	raw, ok := ts["ts_price"].([]any)
	if !ok {
		return nil
	}
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
	// downsample long intraday series so the response stays lean
	const maxPoints = 96
	if len(out) > maxPoints {
		stride := (len(out) + maxPoints - 1) / maxPoints
		sampled := make([][2]float64, 0, maxPoints)
		for i := 0; i < len(out); i += stride {
			sampled = append(sampled, out[i])
		}
		out = sampled
	}
	return out
}

func richStr(v any) string {
	s, _ := v.(string)
	return s
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

func richFloatOrDefault(v any) float64 {
	f, _ := richFloat(v)
	return f
}
