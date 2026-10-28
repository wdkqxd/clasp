package brave

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"

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
		pageURL, _ := m["url"].(string)
		title, _ := m["title"].(string)
		if title == "" {
			title, _ = m["source"].(string)
		}

		thumb := ""
		thumbW, thumbH := 0, 0
		if thumbObj, ok := m["thumbnail"].(map[string]any); ok {
			thumb, _ = thumbObj["src"].(string)
			thumbW = intNum(thumbObj["width"])
			thumbH = intNum(thumbObj["height"])
		}

		imgURL, hosted := originalImage(m)
		if imgURL == "" || thumb == "" {
			continue
		}
		if hosted {
			imgURL = proxyImage(imgURL)
		}

		width, height := 0, 0
		if props, ok := m["properties"].(map[string]any); ok {
			width = intNum(props["width"])
			height = intNum(props["height"])
		}
		if width == 0 {
			width = thumbW
		}
		if height == 0 {
			height = thumbH
		}

		out = append(out, engine.ImageResult{
			Title:     title,
			URL:       pageURL,
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

func originalImage(m map[string]any) (string, bool) {
	props, _ := m["properties"].(map[string]any)

	if u, _ := props["url"].(string); u != "" {
		return u, true
	}
	if u := similarOriginal(m["similar_images_api_url"]); u != "" {
		return u, true
	}
	if thumb, ok := m["thumbnail"].(map[string]any); ok {
		if u, _ := thumb["original"].(string); u != "" {
			return u, true
		}
	}
	u, _ := props["resized"].(string)
	return u, false
}

func proxyImage(imageURL string) string {
	return "https://external-content.duckduckgo.com/iu/?u=" + url.QueryEscape(imageURL)
}

func similarOriginal(v any) string {
	raw, _ := v.(string)
	i := strings.Index(raw, "&url=")
	if i < 0 {
		i = strings.Index(raw, "?url=")
		if i < 0 {
			return ""
		}
	}
	rest := raw[i+len("&url="):]
	if j := strings.IndexByte(rest, '&'); j >= 0 {
		rest = rest[:j]
	}
	decoded, err := url.QueryUnescape(rest)
	if err != nil || !strings.HasPrefix(decoded, "http") {
		return ""
	}
	return decoded
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
