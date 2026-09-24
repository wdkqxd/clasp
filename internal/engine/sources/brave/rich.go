package brave

import (
	"encoding/json"
	"strings"

	"metasearch/internal/engine"
	"metasearch/internal/engine/helpers"
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

var richSubtypeKeys = []struct {
	subtype    string
	payloadKey string
	kind       string
}{
	{"calculator", "calculator", "calculator"},
	{"colorpicker", "colorpicker", "colorpicker"},
	{"weather", "weather", "weather"},
	{"timer", "timer", "timer"},
	{"unitconversion", "unitConversion", "unitConversion"},
	{"timezones", "timezones", "timezones"},
	{"cryptocurrency", "cryptocurrency", "cryptocurrency"},
	{"unixtimestamp", "unixtimestamp", "unixtimestamp"},
	{"currency", "currency", "currency"},
	{"news", "news", "news"},
}

func richAnswerFromMap(m map[string]any) *engine.RichAnswer {
	subtype, _ := m["subtype"].(string)
	subtype = strings.ToLower(subtype)
	if subtype == "stopwatch" {
		return &engine.RichAnswer{Type: "stopwatch", Stopwatch: true}
	}
	for _, sub := range richSubtypeKeys {
		if subtype != "" && subtype != sub.subtype {
			continue
		}
		payload, ok := m[sub.payloadKey].(map[string]any)
		if !ok || len(payload) == 0 {
			continue
		}
		switch sub.kind {
		case "calculator":
			if a := richCalculator(payload); a != nil {
				return a
			}
		case "colorpicker":
			if a := richColorpicker(payload); a != nil {
				return a
			}
		case "weather":
			if a := richWeather(payload); a != nil {
				return a
			}
		case "timer":
			if a := richTimer(payload); a != nil {
				return a
			}
		case "unitConversion":
			if a := richUnitConversion(payload); a != nil {
				return a
			}
		case "timezones":
			if a := richTimezones(payload); a != nil {
				return a
			}
		case "cryptocurrency":
			if a := richCryptocurrency(payload); a != nil {
				return a
			}
		case "unixtimestamp":
			if a := richUnixtimestamp(payload); a != nil {
				return a
			}
		case "currency":
			if a := richCurrency(payload); a != nil {
				return a
			}
		case "news":
			if a := richNews(payload); a != nil {
				return a
			}
		}
		// subtype matched but payload not usable -> stop scanning
		if subtype == sub.subtype {
			break
		}
	}
	return nil
}

func richCalculator(p map[string]any) *engine.RichAnswer {
	expr, _ := p["expression"].(string)
	ans, ok := richFloat(p["answer"])
	if !ok || expr == "" {
		return nil
	}
	return &engine.RichAnswer{Type: "calculator", Calculator: &engine.RichCalculator{Expression: expr, Answer: ans}}
}

func richColorpicker(p map[string]any) *engine.RichAnswer {
	color, _ := p["color"].(string)
	if color == "" {
		return nil
	}
	format, _ := p["format"].(string)
	return &engine.RichAnswer{Type: "colorpicker", Colorpicker: &engine.RichColorpicker{Color: color, Format: format}}
}

func richTimer(p map[string]any) *engine.RichAnswer {
	dur, ok := richFloat(p["duration"])
	if !ok {
		return nil
	}
	start, _ := p["start_on_load"].(bool)
	return &engine.RichAnswer{Type: "timer", Timer: &engine.RichTimer{Duration: int64(dur), StartOnLoad: start}}
}

func richUnitConversion(p map[string]any) *engine.RichAnswer {
	amount, ok := richFloat(p["amount"])
	if !ok {
		return nil
	}
	from, _ := p["from_unit"].(string)
	to, _ := p["to_unit"].(string)
	if from == "" || to == "" {
		return nil
	}
	dim, _ := p["dimensionality"].(string)
	return &engine.RichAnswer{
		Type: "unitConversion",
		UnitConversion: &engine.RichUnitConversion{
			Amount: amount, FromUnit: from, ToUnit: to, Dimensionality: dim,
		},
	}
}

func richTimezones(p map[string]any) *engine.RichAnswer {
	res, ok := p["result"].(map[string]any)
	if !ok {
		return nil
	}
	out := &engine.RichTimezones{Type: richStr(res["type"])}
	raw, ok := res["timezones"].([]any)
	if !ok {
		return nil
	}
	for _, t := range raw {
		tm, ok := t.(map[string]any)
		if !ok {
			continue
		}
		rz := engine.RichTimezone{
			Abbreviation: richStr(tm["abbreviation"]),
			GenericName:  richStr(tm["generic_name"]),
			UTCOffset:    richStr(tm["utc_offset"]),
		}
		if ct, ok := tm["converted_time"].(map[string]any); ok {
			city := richStr(ct["city"])
			if cm, ok := ct["city"].(map[string]any); ok {
				city = richStr(cm["name"])
			}
			rz.ConvertedTime = &engine.RichTimezoneCnv{
				Strftime: richStr(ct["strftime"]),
				Location: richStr(ct["location"]),
				UTCDiff:  richStr(ct["utc_diff"]),
				City:     city,
			}
		}
		out.Timezones = append(out.Timezones, rz)
		if len(out.Timezones) >= 6 {
			break
		}
	}
	if len(out.Timezones) == 0 {
		return nil
	}
	return &engine.RichAnswer{Type: "timezones", Timezones: out}
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

func richUnixtimestamp(p map[string]any) *engine.RichAnswer {
	if _, ok := p["conversion"].(map[string]any); !ok {
		return nil
	}
	return &engine.RichAnswer{Type: "unixtimestamp", Unixtimestamp: &engine.RichUnixtimestamp{Conversion: p["conversion"]}}
}

func richCurrency(p map[string]any) *engine.RichAnswer {
	conv, ok := p["conversion"].(map[string]any)
	if !ok {
		return nil
	}
	out := &engine.RichCurrency{Amount: richFloatOrDefault(conv["amount"])}
	if q, ok := conv["query"].(map[string]any); ok {
		if from, ok := q["from_currency"].(map[string]any); ok {
			out.FromCode = richStr(from["code"])
			out.FromName = richStr(from["full_name"])
		}
		if to, ok := q["to_currency"].(map[string]any); ok {
			out.ToCode = richStr(to["code"])
			out.ToName = richStr(to["full_name"])
		}
		if amt, ok := richFloat(q["amount"]); ok && out.Amount == 0 {
			out.Amount = amt
		}
	}
	out.Result = richFloatOrDefault(conv["result"])
	out.Date = richStr(conv["date"])
	if info, ok := conv["info"].(map[string]any); ok {
		out.ExchangeRate = richFloatOrDefault(info["rate"])
	}
	if ts, ok := p["timeseries"].(map[string]any); ok {
		series := &engine.RichCurrSeries{TimeRange: richStr(ts["time_range"])}
		if raw, ok := ts["ts_exchange_rate"].([]any); ok {
			for _, r := range raw {
				pair, ok := r.([]any)
				if !ok || len(pair) < 2 {
					continue
				}
				series.TsExchange = append(series.TsExchange, [2]any{pair[0], pair[1]})
				if len(series.TsExchange) >= 30 {
					break
				}
			}
		}
		out.Timeseries = series
	}
	if out.FromCode == "" || out.ToCode == "" {
		return nil
	}
	return &engine.RichAnswer{Type: "currency", Currency: out}
}

func richNews(p map[string]any) *engine.RichAnswer {
	raw, ok := p["articles"].([]any)
	if !ok {
		return nil
	}
	out := &engine.RichNews{}
	for _, a := range raw {
		am, ok := a.(map[string]any)
		if !ok {
			continue
		}
		art := engine.RichNewsArticle{
			Title:         richStr(am["title"]),
			URL:           richStr(am["url"]),
			Description:   richStr(am["description"]),
			Img:           richStr(am["img"]),
			PublishTime:   int64(richFloatOrDefault(am["publish_time"])),
			PublisherName: richStr(am["publisher_name"]),
			PageAge:       richStr(am["page_age"]),
		}
		if mu, ok := am["meta_url"].(map[string]any); ok {
			art.MetaURL = &engine.RichMetaURL{Hostname: richStr(mu["hostname"]), Favicon: richStr(mu["favicon"])}
		}
		out.Articles = append(out.Articles, art)
		if len(out.Articles) >= 6 {
			break
		}
	}
	if len(out.Articles) == 0 {
		return nil
	}
	return &engine.RichAnswer{Type: "news", News: out}
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
