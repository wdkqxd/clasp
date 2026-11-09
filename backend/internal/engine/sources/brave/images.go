package brave

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	"clasp/internal/engine"
)

type imageThumbnail struct {
	Src      string
	Original string
	Width    int
	Height   int
}

type imageProperties struct {
	URL     string
	Resized string
	Width   int
	Height  int
}

type imageEntry struct {
	URL                 string
	Title               string
	Source              string
	SimilarImagesAPIURL string
	Thumbnail           imageThumbnail
	Properties          imageProperties
}

func imagePayload(data map[string]any) []imageEntry {
	raw := fields(data).obj("body").obj("response").list("results")
	out := make([]imageEntry, 0, len(raw))
	for _, item := range raw {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		f := fields(m)
		thumb := f.obj("thumbnail")
		props := f.obj("properties")
		out = append(out, imageEntry{
			URL:                 f.str("url"),
			Title:               f.str("title"),
			Source:              f.str("source"),
			SimilarImagesAPIURL: f.str("similar_images_api_url"),
			Thumbnail: imageThumbnail{
				Src:      thumb.str("src"),
				Original: thumb.str("original"),
				Width:    intNum(thumb["width"]),
				Height:   intNum(thumb["height"]),
			},
			Properties: imageProperties{
				URL:     props.str("url"),
				Resized: props.str("resized"),
				Width:   intNum(props["width"]),
				Height:  intNum(props["height"]),
			},
		})
	}
	return out
}

func imageResults(data map[string]any, limit int) ([]engine.ImageResult, error) {
	raw := imagePayload(data)

	out := make([]engine.ImageResult, 0, len(raw))
	for _, e := range raw {
		title := e.Title
		if title == "" {
			title = e.Source
		}

		thumb := e.Thumbnail.Src
		imgURL, hosted := originalImage(e)
		if imgURL == "" || thumb == "" {
			continue
		}
		if hosted {
			imgURL = proxyImage(imgURL)
		}

		width, height := e.Properties.Width, e.Properties.Height
		if width == 0 {
			width = e.Thumbnail.Width
		}
		if height == 0 {
			height = e.Thumbnail.Height
		}

		out = append(out, engine.ImageResult{
			Title:     title,
			URL:       e.URL,
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

func originalImage(e imageEntry) (string, bool) {
	if e.Properties.URL != "" {
		return e.Properties.URL, true
	}
	if u := similarOriginal(e.SimilarImagesAPIURL); u != "" {
		return u, true
	}
	if e.Thumbnail.Original != "" {
		return e.Thumbnail.Original, true
	}
	return e.Properties.Resized, false
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
