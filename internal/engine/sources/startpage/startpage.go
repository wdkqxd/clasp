package startpage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "golang.org/x/image/webp"

	"metasearch/internal/engine"
)

const (
	endpoint = "https://www.startpage.com/sp/search"
	passPath = "/.within.website/x/cmd/anubis/api/pass-challenge"
)

var (
	reAnubisChallenge = regexp.MustCompile(`(?s)<script id="anubis_challenge" type="application/json">\s*(\{.*?\})\s*</script>`)
	reSpTitle         = regexp.MustCompile(`(?s)<a class="result-title result-link[^"]*" href="([^"]+)"[^>]*>\s*<h2 class="wgl-title[^"]*"[^>]*>(.*?)</h2>\s*</a>`)
	reSpDescription   = regexp.MustCompile(`(?s)<p class="description[^"]*"[^>]*>(.*?)</p>`)
	reSpPageForm      = regexp.MustCompile(`(?s)<form action="/sp/search" method="post"[^>]*aria-label="go to page ([0-9]+)"[^>]*>(.*?)</form>`)
	reSpInputTag      = regexp.MustCompile(`(?s)<input\b[^>]*>`)
	reSpAttr          = regexp.MustCompile(`([\w.\-]+)=("([^"]*)"|'([^']*)'|[^\s>]+)`)
)

type Startpage struct {
	client engine.HTTPClient

	mu      sync.Mutex
	cookies map[string]string // cookie jar shared across queries
}

func New(client engine.HTTPClient) *Startpage {
	return &Startpage{client: client, cookies: make(map[string]string)}
}

func (s *Startpage) Name() string { return "startpage" }

func (s *Startpage) ImageSearch(ctx context.Context, query string, limit, offset int) ([]engine.ImageResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("query", query)
	q.Set("cat", "images")
	q.Set("language", "english")
	if offset > 0 {
		q.Set("page", strconv.Itoa(offset+1))
	}
	u.RawQuery = q.Encode()

	body, err := s.get(ctx, u.String())
	if err != nil {
		return nil, err
	}

	if reAnubisChallenge.MatchString(body) {
		if err := s.passChallenge(ctx, body); err != nil {
			return nil, err
		}
		body, err = s.get(ctx, u.String())
		if err != nil {
			return nil, err
		}
		if reAnubisChallenge.MatchString(body) {
			return nil, errors.New("startpage still serving a challenge after proof of work")
		}
	}

	results, err := parseImageResults(body, limit)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	for i := range results {
		wg.Add(1)
		go func(r *engine.ImageResult) {
			defer wg.Done()
			r.Width, r.Height = imageSize(ctx, s.client, r.Thumbnail)
		}(&results[i])
	}
	wg.Wait()

	return results, nil
}

type anubisChallenge struct {
	Rules struct {
		Algorithm  string `json:"algorithm"`
		Difficulty int    `json:"difficulty"`
	} `json:"rules"`
	Challenge struct {
		ID         string `json:"id"`
		RandomData string `json:"randomData"`
		Difficulty int    `json:"difficulty"`
	} `json:"challenge"`
}

func (s *Startpage) Search(ctx context.Context, query string, limit, offset int) ([]engine.Result, error) {
	pageNum := offset + 1

	s.mu.Lock()
	defer s.mu.Unlock()

	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("query", query)
	q.Set("language", "english")
	u.RawQuery = q.Encode()

	body, err := s.get(ctx, u.String())
	if err != nil {
		return nil, err
	}

	if reAnubisChallenge.MatchString(body) {
		if err := s.passChallenge(ctx, body); err != nil {
			return nil, err
		}
		body, err = s.get(ctx, u.String())
		if err != nil {
			return nil, err
		}
		if reAnubisChallenge.MatchString(body) {
			return nil, errors.New("startpage still serving a challenge after proof of work")
		}
	}

	if pageNum == 1 {
		return parseResults(body, limit)
	}

	form, ok := pageForm(body, pageNum)
	if !ok {
		return nil, fmt.Errorf("startpage: no pagination form for page %d", pageNum)
	}
	form.Set("query", query)
	form.Set("page", strconv.Itoa(pageNum))

	body, err = s.post(ctx, endpoint, form)
	if err != nil {
		return nil, err
	}

	return parseResults(body, limit)
}

func (s *Startpage) get(ctx context.Context, u string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}
	setHeaders(req, s.cookies)
	return s.do(ctx, req)
}

func (s *Startpage) post(ctx context.Context, u string, form url.Values) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	setHeaders(req, s.cookies)
	return s.do(ctx, req)
}

func (s *Startpage) do(ctx context.Context, req *http.Request) (string, error) {
	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	storeCookies(s.cookies, resp.Header.Values("Set-Cookie"))

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("startpage returned %s", resp.Status)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, engine.MaxBodyReadLimit))
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (s *Startpage) passChallenge(ctx context.Context, doc string) error {
	m := reAnubisChallenge.FindStringSubmatch(doc)
	if len(m) < 2 {
		return errors.New("startpage: no anubis challenge found")
	}
	var ch anubisChallenge
	if err := json.Unmarshal([]byte(m[1]), &ch); err != nil {
		return err
	}

	diff := ch.Challenge.Difficulty
	if diff < 1 {
		diff = ch.Rules.Difficulty
	}

	start := time.Now()
	nonce, hash := solveAnubis(ch.Challenge.RandomData, diff)
	if hash == "" {
		return errors.New("startpage: failed to solve proof of work")
	}

	v := url.Values{}
	v.Set("id", ch.Challenge.ID)
	v.Set("response", hash)
	v.Set("nonce", strconv.FormatInt(nonce, 10))
	v.Set("redir", endpoint)
	v.Set("elapsedTime", strconv.FormatInt(time.Since(start).Milliseconds(), 10))
	passURL := "https://www.startpage.com" + passPath + "?" + v.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, passURL, nil)
	if err != nil {
		return err
	}
	setHeaders(req, s.cookies)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	storeCookies(s.cookies, resp.Header.Values("Set-Cookie"))

	if s.cookies["spchal-auth"] == "" {
		return errors.New("startpage proof of work was rejected (no spchal-auth cookie)")
	}
	return nil
}

func solveAnubis(randomData string, difficulty int) (int64, string) {
	prefix := strings.Repeat("0", difficulty)
	for nonce := int64(0); nonce < 1<<36; nonce++ {
		sum := sha256.Sum256([]byte(randomData + strconv.FormatInt(nonce, 10)))
		digest := hex.EncodeToString(sum[:])
		if strings.HasPrefix(digest, prefix) {
			return nonce, digest
		}
	}
	return 0, ""
}

func setHeaders(req *http.Request, cookies map[string]string) {
	req.Header.Set("User-Agent", engine.BrowserUserAgent)
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("accept-encoding", "gzip")
	req.Header.Set("sec-ch-ua", `"Not(A:Brand";v="99", "Chromium";v="150", "Google Chrome";v="150"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "none")
	if len(cookies) > 0 {
		parts := make([]string, 0, len(cookies))
		for k, v := range cookies {
			parts = append(parts, k+"="+v)
		}
		req.Header.Set("Cookie", strings.Join(parts, "; "))
	}
}

var reSetCookieName = regexp.MustCompile(`^([^=;]+)=([^;]*)`)

func storeCookies(jar map[string]string, lines []string) {
	for _, line := range lines {
		m := reSetCookieName.FindStringSubmatch(line)
		if len(m) < 3 {
			continue
		}
		name := strings.TrimSpace(m[1])
		value := strings.TrimSpace(m[2])
		if value == "" {
			delete(jar, name)
		} else {
			jar[name] = value
		}
	}
}

func pageForm(doc string, page int) (url.Values, bool) {
	for _, m := range reSpPageForm.FindAllStringSubmatch(doc, -1) {
		p, err := strconv.Atoi(m[1])
		if err != nil || p != page {
			continue
		}
		form := url.Values{}
		for _, tag := range reSpInputTag.FindAllString(m[2], -1) {
			attrs := make(map[string]string)
			for _, a := range reSpAttr.FindAllStringSubmatch(tag, -1) {
				v := a[3]
				if v == "" {
					v = a[4]
				}
				attrs[a[1]] = v
			}
			if name, ok := attrs["name"]; ok {
				form.Set(name, attrs["value"])
			}
		}
		return form, true
	}
	return nil, false
}

func parseResults(doc string, limit int) ([]engine.Result, error) {
	idxs := reSpTitle.FindAllStringSubmatchIndex(doc, -1)
	if len(idxs) == 0 {
		return nil, errors.New("no startpage results found in page")
	}

	out := make([]engine.Result, 0, len(idxs))
	for i, m := range idxs {
		resultURL := doc[m[2]:m[3]]
		if !strings.HasPrefix(resultURL, "http://") && !strings.HasPrefix(resultURL, "https://") {
			continue
		}
		r := engine.Result{
			Title:  cleanText(doc[m[4]:m[5]]),
			URL:    resultURL,
			Engine: "startpage",
		}

		segStart := m[1]
		segEnd := len(doc)
		if i+1 < len(idxs) {
			segEnd = idxs[i+1][0]
		}
		if d := reSpDescription.FindStringSubmatch(doc[segStart:segEnd]); len(d) > 1 {
			r.Snippet = cleanText(d[1])
		}

		out = append(out, r)
		if len(out) >= limit {
			break
		}
	}
	if len(out) == 0 {
		return nil, errors.New("no usable startpage results found")
	}
	return out, nil
}

var (
	reSpTags           = regexp.MustCompile(`(?s)<[^>]+>`)
	reSpImageContainer = regexp.MustCompile(`(?s)<div[^>]*class="image-container grid-image-\d+[^"]*"[^>]*data-thumbnail-url="([^"]*)"[^>]*aria-label="click to expand image: ([^"]*)"`)
	reSpImageSource    = regexp.MustCompile(`<span[^>]*>([^<]+)</span>`)
)

func cleanText(raw string) string {
	raw = reSpTags.ReplaceAllString(raw, "")
	return html.UnescapeString(strings.TrimSpace(raw))
}

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
