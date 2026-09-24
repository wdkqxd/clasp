package brave

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"metasearch/internal/engine"
)

const (
	webEndpoint = "https://search.brave.com/search"
	imageURI    = "https://search.brave.com/images"
)

type Brave struct {
	client engine.HTTPClient

	mu          sync.Mutex
	infobox     *engine.Infobox
	discussions []engine.Discussion
	rich        []engine.RichAnswer
}

func New(client engine.HTTPClient) *Brave {
	return &Brave{client: client}
}

func (b *Brave) Name() string { return "brave" }

func (b *Brave) Infobox() *engine.Infobox {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.infobox
}

func (b *Brave) Discussions() []engine.Discussion {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.discussions
}

func (b *Brave) RichAnswers() []engine.RichAnswer {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.rich
}

func (b *Brave) ImageSearch(ctx context.Context, query string, limit, offset int) ([]engine.ImageResult, error) {
	doc, err := b.fetchDoc(ctx, "brave images", imageURI, query, offset)
	if err != nil {
		return nil, err
	}
	data, err := parseSerpData(doc)
	if err != nil {
		return nil, err
	}
	return imageResults(data, limit)
}

func (b *Brave) Search(ctx context.Context, query string, limit, offset int) ([]engine.Result, error) {
	b.mu.Lock()
	b.infobox = nil
	b.discussions = nil
	b.rich = nil
	b.mu.Unlock()

	doc, err := b.fetchDoc(ctx, "brave web", webEndpoint, query, offset)
	if err != nil {
		return nil, err
	}

	data, serpErr := parseSerpData(doc)
	b.mu.Lock()
	b.infobox = parseInfobox(doc)
	b.discussions = parseDiscussions(doc)
	b.rich = parseRich(doc)
	b.mu.Unlock()

	if serpErr == nil {
		return webResults(data, limit)
	}
	return ssrResults(doc, limit)
}

func (b *Brave) fetchDoc(ctx context.Context, label, rawURL, query string, offset int) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("q", query)
	q.Set("source", "web")
	q.Set("spellcheck", "0")
	q.Set("search_lang", "en")
	q.Set("country", "ALL")
	if offset > 0 {
		q.Set("offset", strconv.Itoa(offset))
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", err
	}
	setHeaders(req)

	resp, err := b.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s returned %s", label, resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, engine.MaxBodyReadLimit))
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func setHeaders(req *http.Request) {
	req.Header.Set("User-Agent", engine.BrowserUserAgent)
	req.Header.Set("Cookie", "theme=dark; country=all; useLocation=0")
	req.Header.Set("sec-ch-ua", `"Chromium";v="133", "Not(A:Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("sec-gpc", "1")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("referer", "https://search.brave.com/")
}
