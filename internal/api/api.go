package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"clasp/internal/aggregator"
	"clasp/internal/bangs"
	"clasp/internal/suggest"
	"clasp/internal/wikipedia"
)

type Searcher interface {
	Search(ctx context.Context, query string, limit, page int, engines []string) aggregator.Response
	ImageSearch(ctx context.Context, query string, limit, page int, engines []string) aggregator.ImageResponse
}

type Server struct {
	search    Searcher
	bangs     *bangs.Store
	wikipedia *wikipedia.Client
	suggest   *suggest.Client
	start     time.Time
	corsOrig  string
}

func New(search Searcher, b *bangs.Store, wiki *wikipedia.Client, suggestions *suggest.Client) http.Handler {
	s := &Server{
		search:    search,
		bangs:     b,
		wikipedia: wiki,
		suggest:   suggestions,
		start:     time.Now(),
		corsOrig:  os.Getenv("CORS_ORIGIN"),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", s.handleHealth)
	mux.HandleFunc("GET /api/wikipedia", s.handleWikipedia)
	mux.HandleFunc("GET /api/search", s.handleSearch)
	mux.HandleFunc("GET /api/images", s.handleImages)
	mux.HandleFunc("GET /api/bangs", s.handleBangs)
	mux.HandleFunc("GET /api/suggest", s.handleSuggest)
	return s.withCORS(mux)
}

func (s *Server) withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := s.corsOrig
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
		"uptime": time.Since(s.start).Truncate(time.Second).String(),
	})
}

func (s *Server) handleBangs(w http.ResponseWriter, r *http.Request) {
	s.bangs.ServeJSON(w)
}

type searchParams struct {
	query   string
	limit   int
	page    int
	engines []string
}

func (s *Server) parseSearchParams(w http.ResponseWriter, r *http.Request, defaultLimit, maxLimit int) (searchParams, bool) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "missing required query parameter 'q'",
		})
		return searchParams{}, false
	}

	limit := defaultLimit
	if raw := r.URL.Query().Get("limit"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"error": "invalid 'limit': must be an integer",
			})
			return searchParams{}, false
		}
		limit = n
	}
	if limit < 1 {
		limit = 1
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	page := 0
	if raw := r.URL.Query().Get("offset"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 0 {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"error": "invalid 'offset': must be a non-negative integer",
			})
			return searchParams{}, false
		}
		page = n
	} else if raw := r.URL.Query().Get("page"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 0 {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"error": "invalid 'page': must be a non-negative integer",
			})
			return searchParams{}, false
		}
		page = n
	}

	var engines []string
	if raw := strings.TrimSpace(r.URL.Query().Get("engines")); raw != "" {
		for _, name := range strings.Split(raw, ",") {
			if name = strings.TrimSpace(name); name != "" {
				engines = append(engines, name)
			}
		}
	}

	return searchParams{query: query, limit: limit, page: page, engines: engines}, true
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	params, ok := s.parseSearchParams(w, r, 10, 50)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, s.search.Search(r.Context(), params.query, params.limit, params.page, params.engines))
}

func (s *Server) handleImages(w http.ResponseWriter, r *http.Request) {
	params, ok := s.parseSearchParams(w, r, 20, 200)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, s.search.ImageSearch(r.Context(), params.query, params.limit, params.page, params.engines))
}

func (s *Server) handleSuggest(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "missing required query parameter 'q'",
		})
		return
	}

	engine := suggest.Normalize(r.URL.Query().Get("engine"))
	if engine == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "invalid 'engine': must be one of " + strings.Join(suggest.Engines(), ", "),
		})
		return
	}

	resp, err := s.suggest.Suggest(r.Context(), query, engine)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{
			"engine": engine,
			"error":  err.Error(),
		})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
