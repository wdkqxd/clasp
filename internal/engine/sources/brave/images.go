package brave

import (
	"encoding/json"
	"errors"

	"clasp/internal/engine"
)

func imageResults(data map[string]any, limit int) ([]engine.ImageResult, error) {
	body, _ := data["body"].(map[string]any)
	response, _ := body["response"].(map[string]any)
	raw, _ := response["results"].([]any)

	out := make([]engine.ImageResult, 0, len(raw))
	for _, item := range raw {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		title, _ := m["title"].(string)
		sourceURL, _ := m["url"].(string)

		// thumbnail contains the small image
		thumb := ""
		thumbW, thumbH := 0, 0
		if thumbObj, ok := m["thumbnail"].(map[string]any); ok {
			thumb, _ = thumbObj["src"].(string)
			thumbW = intNum(thumbObj["width"])
			thumbH = intNum(thumbObj["height"])
		}

		// properties contains the full-size image
		imgURL := sourceURL
		width, height := thumbW, thumbH
		if props, ok := m["properties"].(map[string]any); ok {
			if u, ok := props["url"].(string); ok && u != "" {
				imgURL = u
			}
			width = intNum(props["width"])
			height = intNum(props["height"])
		}

		if imgURL == "" && thumb == "" {
			continue
		}

		out = append(out, engine.ImageResult{
			Title:     title,
			URL:       sourceURL,
			SourceURL: imgURL,
			Thumbnail: thumb,
			Width:     width,
			Height:    height,
			Engine:    "brave",
		})
		if len(out) >= limit {
			break
		}
	}
	if len(out) == 0 {
		return nil, errors.New("no image results found in brave payload")
	}
	return out, nil
}

func intNum(v any) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case int64:
		return int(n)
	case json.Number:
		i, _ := n.Int64()
		return int(i)
	}
	return 0
}
