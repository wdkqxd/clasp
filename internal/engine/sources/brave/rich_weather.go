package brave

import (
	"metasearch/internal/engine"
)

func richWeather(p map[string]any) *engine.RichAnswer {
	out := &engine.RichWeather{}
	if loc, ok := p["location"].(map[string]any); ok {
		out.Location = &engine.RichWeatherLocation{
			Name:    richStr(loc["name"]),
			Country: richStr(loc["country"]),
			State:   richStr(loc["state"]),
		}
	}
	if cw, ok := p["current_weather"].(map[string]any); ok {
		cur := &engine.RichWeatherCurrent{
			Temp:      richFloatOrDefault(cw["temp"]),
			FeelsLike: richFloatOrDefault(cw["feels_like"]),
			Humidity:  richFloatOrDefault(cw["humidity"]),
			UVI:       richFloatOrDefault(cw["uvi"]),
			Pressure:  richFloatOrDefault(cw["pressure"]),
		}
		if wind, ok := cw["wind"].(map[string]any); ok {
			cur.WindSpeed = richFloatOrDefault(wind["speed"])
		}
		if wth, ok := cw["weather"].(map[string]any); ok {
			cur.Weather = &engine.RichWeatherCond{
				Main:        richStr(wth["main"]),
				Description: richStr(wth["description"]),
				Icon:        richStr(wth["icon"]),
			}
		}
		out.CurrentWeather = cur
	}
	if days, ok := p["daily"].([]any); ok {
		for _, d := range days {
			dm, ok := d.(map[string]any)
			if !ok {
				continue
			}
			day := engine.RichWeatherDay{
				Ts:       int64(richFloatOrDefault(dm["ts"])),
				DateI18n: richStr(dm["date_i18n"]),
			}
			if t, ok := dm["temperature"].(map[string]any); ok {
				day.Temperature = &engine.RichWeatherTemp{
					Max: richFloatOrDefault(t["max"]),
					Min: richFloatOrDefault(t["min"]),
					Day: richFloatOrDefault(t["day"]),
				}
			}
			if wth, ok := dm["weather"].(map[string]any); ok {
				day.Weather = &engine.RichWeatherCond{
					Main:        richStr(wth["main"]),
					Description: richStr(wth["description"]),
					Icon:        richStr(wth["icon"]),
				}
			}
			out.Daily = append(out.Daily, day)
			if len(out.Daily) >= 5 {
				break
			}
		}
	}
	if alerts, ok := p["alerts"].([]any); ok {
		for _, a := range alerts {
			am, ok := a.(map[string]any)
			if !ok {
				continue
			}
			out.Alerts = append(out.Alerts, engine.RichWeatherAlert{
				Event:             richStr(am["event"]),
				Description:       richStr(am["description"]),
				StartRelativeI18n: richStr(am["start_relative_i18n"]),
			})
			if len(out.Alerts) >= 2 {
				break
			}
		}
	}
	if out.CurrentWeather == nil && len(out.Daily) == 0 {
		return nil
	}
	return &engine.RichAnswer{Type: "weather", Weather: out}
}
