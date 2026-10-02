package aggregator

import (
	"context"
	"fmt"
	"net/url"
	"slices"
	"sort"
	"strings"
	"time"

	"metasearch/internal/engine"
)

type Status struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
	TookMS int64  `json:"took_ms"`
	Error  string `json:"error,omitempty"`
}

type Response struct {
	Query       string              `json:"query"`
	TookMS      int64               `json:"took_ms"`
	Results     []engine.Result     `json:"results"`
	Engines     map[string]Status   `json:"engines"`
	HasMore     bool                `json:"has_more"`
	Infobox     *engine.Infobox     `json:"infobox,omitempty"`
	Discussions []engine.Discussion `json:"discussions,omitempty"`
	Rich        []engine.RichAnswer `json:"rich,omitempty"`
}

type ImageResponse struct {
	Query   string               `json:"query"`
	TookMS  int64                `json:"took_ms"`
	Results []engine.ImageResult `json:"results"`
	Engines map[string]Status    `json:"engines"`
	HasMore bool                 `json:"has_more"`
}

type Aggregator struct {
	engines          []engine.Engine
	totalTimeout     time.Duration
	perEngineTimeout time.Duration
	fastTimeout      time.Duration
}

func New(engines []engine.Engine) *Aggregator {
	return &Aggregator{
		engines:          engines,
		totalTimeout:     10 * time.Second,
		perEngineTimeout: 5 * time.Second,
		fastTimeout:      3 * time.Second,
	}
}

// engineOutcome carries an engine result together with the engine it came from.
type engineOutcome[T any] struct {
	name string
	val  T
}

// collectOutcomes waits for one outcome per engine, but gives up on stragglers
// once timeout elapses instead of waiting for the slowest source. Engines that
// did not report in time get fallback(name), so callers can mark them failed.
func collectOutcomes[T any](names []string, ch <-chan engineOutcome[T], timeout time.Duration, fallback func(string) T) map[string]T {
	pending := make(map[string]struct{}, len(names))
	for _, n := range names {
		pending[n] = struct{}{}
	}

	out := make(map[string]T, len(names))
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for len(pending) > 0 {
		select {
		case o := <-ch:
			out[o.name] = o.val
			delete(pending, o.name)
		case <-timer.C:
			for n := range pending {
				out[n] = fallback(n)
			}
			return out
		}
	}
	return out
}

const (
	rrfK              = 10
	maxCandidates     = 60
	imageCandidates   = 200
	navigationalBoost = 4.0
)

func selectEngines(all []engine.Engine, enabled []string) []engine.Engine {
	if len(enabled) == 0 {
		return all
	}
	selected := make(map[string]bool, len(enabled))
	for _, name := range enabled {
		name = strings.ToLower(strings.TrimSpace(name))
		if name == "duckduckgo" {
			name = "ddg"
		}
		selected[name] = true
	}
	out := make([]engine.Engine, 0, len(all))
	for _, e := range all {
		if selected[strings.ToLower(e.Name())] {
			out = append(out, e)
		}
	}
	return out
}

func (a *Aggregator) Search(ctx context.Context, query string, limit, page int, enabled []string) Response {
	start := time.Now()
	ctx, cancel := context.WithTimeout(ctx, a.totalTimeout)
	defer cancel()

	if limit < 1 {
		limit = 1
	}
	if limit > maxCandidates {
		limit = maxCandidates
	}
	if page < 0 {
		page = 0
	}

	activeEngines := selectEngines(a.engines, enabled)

	resp := Response{
		Query:   query,
		Engines: make(map[string]Status, len(activeEngines)),
	}

	type outcome struct {
		res  engine.SearchResult
		err  error
		took time.Duration
	}

	names := make([]string, 0, len(activeEngines))
	ch := make(chan engineOutcome[outcome], len(activeEngines))
	for _, e := range activeEngines {
		names = append(names, e.Name())
		go func(e engine.Engine) {
			ectx, cancel := context.WithTimeout(ctx, a.perEngineTimeout)
			defer cancel()
			start := time.Now()
			res, err := e.Search(ectx, query, maxCandidates, page)
			ch <- engineOutcome[outcome]{
				name: e.Name(),
				val:  outcome{res: res, err: err, took: time.Since(start)},
			}
		}(e)
	}

	timeoutErr := fmt.Errorf("timeout after %s", a.fastTimeout)
	byName := collectOutcomes(names, ch, a.fastTimeout, func(string) outcome {
		return outcome{err: timeoutErr, took: a.fastTimeout}
	})

	seen := make(map[string]int)

	type scored struct {
		res   engine.Result
		score float64
		nav   bool
	}
	scoredResults := make([]scored, 0)

	navTerm := navigationalTerm(query)
	seenNav := make(map[string]int)

	for _, e := range activeEngines {
		o := byName[e.Name()]
		st := Status{TookMS: o.took.Milliseconds()}
		if o.err != nil {
			st.Status = "error"
			st.Error = o.err.Error()
		} else {
			st.Status = "ok"
			st.Count = len(o.res.Results)
			for j, r := range o.res.Results {
				key := normalizeURL(r.URL)
				if key == "" {
					continue
				}
				if r.Engine != "" && len(r.Engines) == 0 {
					r.Engines = []string{r.Engine}
				}
				contrib := 1.0 / (float64(j+1) + rrfK)
				navRoot := navTerm != "" && navRootMatch(r.URL, navTerm)
				navKey := key
				if navRoot {
					navKey = baseKey(key)
					if idx, dup := seenNav[navKey]; dup {
						scoredResults[idx].score += contrib
						mergeTwoResults(&scoredResults[idx].res, r)
						continue
					}
				}
				if idx, dup := seen[key]; dup {
					scoredResults[idx].score += contrib
					mergeTwoResults(&scoredResults[idx].res, r)
					if navRoot {
						scoredResults[idx].nav = true
						seenNav[navKey] = idx
					}
					continue
				}
				seen[key] = len(scoredResults)
				if navRoot {
					seenNav[navKey] = len(scoredResults)
				}
				scoredResults = append(scoredResults, scored{res: r, score: contrib, nav: navRoot})
			}
		}
		resp.Engines[e.Name()] = st
	}

	sort.SliceStable(scoredResults, func(x, y int) bool {
		sx, sy := scoredResults[x].score, scoredResults[y].score
		if scoredResults[x].nav {
			sx *= navigationalBoost
		}
		if scoredResults[y].nav {
			sy *= navigationalBoost
		}
		return sx > sy
	})

	hasMore := len(scoredResults) > limit
	if !hasMore {
		for _, o := range byName {
			if o.err == nil && len(o.res.Results) > 0 {
				hasMore = true
				break
			}
		}
	}
	resp.HasMore = hasMore

	for _, e := range activeEngines {
		o := byName[e.Name()]
		if o.err != nil {
			continue
		}
		if resp.Infobox == nil && o.res.Infobox != nil {
			resp.Infobox = o.res.Infobox
		}
		if resp.Discussions == nil && len(o.res.Discussions) > 0 {
			resp.Discussions = o.res.Discussions
		}
		if resp.Rich == nil && len(o.res.Rich) > 0 {
			resp.Rich = o.res.Rich
		}
	}

	if len(scoredResults) > limit {
		scoredResults = scoredResults[:limit]
	}

	resp.Results = make([]engine.Result, 0, len(scoredResults))
	for _, s := range scoredResults {
		resp.Results = append(resp.Results, s.res)
	}

	resp.TookMS = time.Since(start).Milliseconds()
	return resp
}

var trackingParams = []string{
	"utm_source", "utm_medium", "utm_campaign", "utm_term", "utm_content",
	"fbclid", "gclid", "msclkid", "yclid", "ref", "ref_src", "ref_url", "spm",
}

func normalizeURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	if u.Scheme != "" && u.Scheme != "http" && u.Scheme != "https" {
		return ""
	}
	host := strings.TrimPrefix(strings.ToLower(u.Host), "www.")
	if host == "" {
		return ""
	}
	if port := u.Port(); port != "" {
		if (u.Scheme == "http" && port == "80") || (u.Scheme == "https" && port == "443") {
			host = strings.TrimSuffix(host, ":"+port)
		}
	}
	q := u.Query()
	for _, k := range trackingParams {
		q.Del(k)
	}
	u.Scheme = ""
	u.Host = ""
	u.Fragment = ""
	u.RawFragment = ""
	u.Path = strings.ToLower(strings.TrimSuffix(u.Path, "/"))
	u.RawQuery = q.Encode()
	return host + u.String()
}

func navigationalTerm(query string) string {
	q := strings.ToLower(strings.TrimSpace(query))
	if len(q) < 2 || len(q) > 40 {
		return ""
	}
	for _, r := range q {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '.', r == '-':
		default:
			return ""
		}
	}
	if i := strings.LastIndex(q, "."); i > 0 {
		tld := q[i+1:]
		if len(tld) >= 2 && len(tld) <= 4 && isAlpha(tld) {
			q = q[:i]
		}
	}
	if q == "" {
		return ""
	}
	return q
}

func isAlpha(s string) bool {
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) {
			return false
		}
	}
	return true
}

func navRootMatch(raw string, term string) bool {
	u, err := url.Parse(strings.ToLower(raw))
	if err != nil {
		return false
	}
	host := strings.Split(u.Host, ":")[0]
	if host == "" {
		return false
	}
	labels := strings.Split(host, ".")
	if len(labels) < 2 || labels[len(labels)-2] != term {
		return false
	}
	if sub := strings.Join(labels[:len(labels)-2], "."); sub != "" && sub != "www" {
		return false
	}
	p := u.Path
	return p == "" || p == "/"
}

func baseKey(key string) string {
	if i := strings.Index(key, "?"); i >= 0 {
		key = key[:i]
	}
	return strings.TrimSuffix(key, "/")
}

func mergeTwoResults(existing *engine.Result, incoming engine.Result) {
	if len(incoming.Title) > len(existing.Title) {
		existing.Title = incoming.Title
	}
	if len(incoming.Snippet) > len(existing.Snippet) {
		existing.Snippet = incoming.Snippet
	}
	// prefer HTTPS when available
	if strings.HasPrefix(existing.URL, "http://") && strings.HasPrefix(incoming.URL, "https://") {
		existing.URL = incoming.URL
	}
	if incoming.Sitelinks != nil && len(incoming.Sitelinks) > len(existing.Sitelinks) {
		existing.Sitelinks = incoming.Sitelinks
	}
	if incoming.Engine != "" && !slices.Contains(existing.Engines, incoming.Engine) {
		existing.Engines = append(existing.Engines, incoming.Engine)
	}
}

func mergeTwoImageResults(existing *engine.ImageResult, incoming engine.ImageResult) {
	if len(incoming.Title) > len(existing.Title) {
		existing.Title = incoming.Title
	}
	if existing.URL == "" {
		existing.URL = incoming.URL
	}
	if existing.SourceURL == "" {
		existing.SourceURL = incoming.SourceURL
	}
	if existing.Thumbnail == "" {
		existing.Thumbnail = incoming.Thumbnail
	}
	if incoming.Engine != "" && !slices.Contains(existing.Engines, incoming.Engine) {
		existing.Engines = append(existing.Engines, incoming.Engine)
	}
}

func (a *Aggregator) ImageSearch(ctx context.Context, query string, limit, page int, enabled []string) ImageResponse {
	start := time.Now()
	ctx, cancel := context.WithTimeout(ctx, a.totalTimeout)
	defer cancel()

	if limit < 1 {
		limit = 1
	}
	if limit > 200 {
		limit = 200
	}
	if page < 0 {
		page = 0
	}

	type imgOutcome struct {
		res  []engine.ImageResult
		err  error
		took time.Duration
	}

	var activeSearchers []struct {
		name string
		s    engine.ImageSearcher
	}
	for _, e := range selectEngines(a.engines, enabled) {
		if is, ok := e.(engine.ImageSearcher); ok {
			activeSearchers = append(activeSearchers, struct {
				name string
				s    engine.ImageSearcher
			}{name: e.Name(), s: is})
		}
	}

	resp := ImageResponse{
		Query:   query,
		Engines: make(map[string]Status, len(activeSearchers)),
	}

	names := make([]string, 0, len(activeSearchers))
	ch := make(chan engineOutcome[imgOutcome], len(activeSearchers))
	for _, as := range activeSearchers {
		names = append(names, as.name)
		go func(name string, s engine.ImageSearcher) {
			ectx, cancel := context.WithTimeout(ctx, a.perEngineTimeout)
			defer cancel()
			start := time.Now()
			res, err := s.ImageSearch(ectx, query, imageCandidates, page)
			ch <- engineOutcome[imgOutcome]{
				name: name,
				val:  imgOutcome{res: res, err: err, took: time.Since(start)},
			}
		}(as.name, as.s)
	}

	timeoutErr := fmt.Errorf("timeout after %s", a.fastTimeout)
	byName := collectOutcomes(names, ch, a.fastTimeout, func(string) imgOutcome {
		return imgOutcome{err: timeoutErr, took: a.fastTimeout}
	})

	seen := make(map[string]int)

	// collect engine statuses
	for _, as := range activeSearchers {
		o := byName[as.name]
		st := Status{TookMS: o.took.Milliseconds()}
		if o.err != nil {
			st.Status = "error"
			st.Error = o.err.Error()
		} else {
			st.Status = "ok"
			st.Count = len(o.res)
		}
		resp.Engines[as.name] = st
	}

	// interleave results round-robin
	type enginePos struct {
		name string
		pos  int
	}
	var positions []enginePos
	for _, as := range activeSearchers {
		o := byName[as.name]
		if o.err == nil && len(o.res) > 0 {
			positions = append(positions, enginePos{name: as.name, pos: 0})
		}
	}
	for len(positions) > 0 && len(resp.Results) < limit {
		var next []enginePos
		for _, ep := range positions {
			o := byName[ep.name]
			for ep.pos < len(o.res) {
				r := o.res[ep.pos]
				ep.pos++
				key := r.Thumbnail
				if key == "" {
					key = r.URL
				}
				if key == "" {
					continue
				}
				if r.Engine != "" && len(r.Engines) == 0 {
					r.Engines = []string{r.Engine}
				}
				if idx, dup := seen[key]; dup {
					mergeTwoImageResults(&resp.Results[idx], r)
					continue
				}
				seen[key] = len(resp.Results)
				resp.Results = append(resp.Results, r)
				break
			}
			if ep.pos < len(o.res) {
				next = append(next, ep)
			}
		}
		positions = next
	}

	hasMore := false
	for _, as := range activeSearchers {
		o := byName[as.name]
		if o.err == nil && len(o.res) >= limit {
			hasMore = true
			break
		}
	}
	resp.HasMore = hasMore

	resp.TookMS = time.Since(start).Milliseconds()
	return resp
}
