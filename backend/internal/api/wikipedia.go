package api

import (
	"net/http"
	"strings"
)

func (s *Server) handleWikipedia(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing required query parameter 'q'"})
		return
	}

	language := r.URL.Query().Get("lang")
	article, err := s.wikipedia.Search(r.Context(), query, language)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "wikipedia unavailable"})
		return
	}
	writeJSON(w, http.StatusOK, article)
}
