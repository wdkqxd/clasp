package startpage

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"sync"

	"clasp/internal/engine"
)

const (
	endpoint = "https://www.startpage.com/sp/search"
	passPath = "/.within.website/x/cmd/anubis/api/pass-challenge"
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
	body, err = s.passIfChallenged(ctx, body, u.String())
	if err != nil {
		return nil, err
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

func (s *Startpage) Search(ctx context.Context, query string, limit, offset int) (engine.SearchResult, error) {
	pageNum := offset + 1

	s.mu.Lock()
	defer s.mu.Unlock()

	u, err := url.Parse(endpoint)
	if err != nil {
		return engine.SearchResult{}, err
	}
	q := u.Query()
	q.Set("query", query)
	q.Set("language", "english")
	u.RawQuery = q.Encode()

	body, err := s.get(ctx, u.String())
	if err != nil {
		return engine.SearchResult{}, err
	}
	body, err = s.passIfChallenged(ctx, body, u.String())
	if err != nil {
		return engine.SearchResult{}, err
	}

	if pageNum == 1 {
		results, err := parseResults(body, limit)
		return engine.SearchResult{Results: results}, err
	}

	form, ok := pageForm(body, pageNum)
	if !ok {
		return engine.SearchResult{}, errors.New("startpage: no pagination form for page " + strconv.Itoa(pageNum))
	}
	form.Set("query", query)
	form.Set("page", strconv.Itoa(pageNum))

	body, err = s.post(ctx, endpoint, form)
	if err != nil {
		return engine.SearchResult{}, err
	}

	results, err := parseResults(body, limit)
	return engine.SearchResult{Results: results}, err
}

// passIfChallenged solves the anubis proof of work when startpage serves a
// challenge, then refetches the page.
func (s *Startpage) passIfChallenged(ctx context.Context, body, rawURL string) (string, error) {
	if !hasChallenge(body) {
		return body, nil
	}
	if err := s.passChallenge(ctx, body); err != nil {
		return "", err
	}
	body, err := s.get(ctx, rawURL)
	if err != nil {
		return "", err
	}
	if hasChallenge(body) {
		return "", errors.New("startpage still serving a challenge after proof of work")
	}
	return body, nil
}
