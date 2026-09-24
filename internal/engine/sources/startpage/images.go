package startpage

import (
	"context"
	"errors"
	"html"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"regexp"
	"strings"

	_ "golang.org/x/image/webp"

	"metasearch/internal/engine"
)

var (
	reSpImageContainer = regexp.MustCompile(`(?s)<div[^>]*class="image-container grid-image-\d+[^"]*"[^>]*data-thumbnail-url="([^"]*)"[^>]*aria-label="click to expand image: ([^"]*)"`)
	reSpImageSource    = regexp.MustCompile(`<span[^>]*>([^<]+)</span>`)
)

func imageSize(ctx context.Context, client engine.HTTPClient, imageURL string) (int, int) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		return 0, 0
	}
	req.Header.Set("User-Agent", engine.BrowserUserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0
	}

	// follow redirects manually (client has WithNotFollowRedirects)
	for i := 0; i < 5 && resp.StatusCode >= 300 && resp.StatusCode < 400; i++ {
		loc := resp.Header.Get("Location")
		resp.Body.Close()
		if loc == "" {
			return 0, 0
		}
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, loc, nil)
		if err != nil {
			return 0, 0
		}
		req.Header.Set("User-Agent", engine.BrowserUserAgent)
		resp, err = client.Do(req)
		if err != nil {
			return 0, 0
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0
	}
	cfg, _, err := image.DecodeConfig(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return 0, 0
	}
	return cfg.Width, cfg.Height
}

func parseImageResults(doc string, limit int) ([]engine.ImageResult, error) {
	out := make([]engine.ImageResult, 0)

	matches := reSpImageContainer.FindAllStringSubmatchIndex(doc, -1)
	for _, loc := range matches {
		thumb := html.UnescapeString(doc[loc[2]:loc[3]])
		title := html.UnescapeString(doc[loc[4]:loc[5]])
		if thumb == "" {
			continue
		}
		if strings.HasPrefix(thumb, "/") {
			thumb = "https://www.startpage.com" + thumb
		}

		source := ""
		after := doc[loc[1]:]
		if sm := reSpImageSource.FindStringSubmatch(after); sm != nil {
			source = strings.TrimSpace(sm[1])
		}

		sourceURL := ""
		if source != "" {
			sourceURL = "https://" + source
		}

		out = append(out, engine.ImageResult{
			Title:     title,
			URL:       sourceURL,
			SourceURL: thumb,
			Thumbnail: thumb,
			Engine:    "startpage",
		})
		if len(out) >= limit {
			break
		}
	}

	if len(out) == 0 {
		return nil, errors.New("no image results found in startpage")
	}
	return out, nil
}