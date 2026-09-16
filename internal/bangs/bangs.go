package bangs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

const sourceURL = "https://raw.githubusercontent.com/imputnet/helium-services/main/svc/bangs/bangs.json"

type bangEntry struct {
	URL string   `json:"u"`
	Ts  []string `json:"ts"`
}

type Store struct {
	mu     sync.RWMutex
	byKey  map[string]string
	client *http.Client
}

func New(ctx context.Context) *Store {
	s := &Store{client: &http.Client{Timeout: 20 * time.Second}}
	s.Refresh(ctx)
	return s
}

func (s *Store) Refresh(ctx context.Context) {
	byKey, err := s.fetch(ctx)
	if err != nil {
		log.Printf("bangs: upstream fetch failed: %v", err)
		if s.byKey == nil {
			s.byKey = make(map[string]string)
		}
		return
	}
	s.mu.Lock()
	s.byKey = byKey
	s.mu.Unlock()
}

func (s *Store) ServeJSON(w http.ResponseWriter) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(s.byKey)
}

var trailingComma = regexp.MustCompile(`,\s*([\]}])`)

func (s *Store) fetch(ctx context.Context) (map[string]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	cleaned := regexp.MustCompile(`(?m)^\s*//.*$`).ReplaceAll(raw, []byte{})

	var entries []bangEntry
	if err := json.Unmarshal(trailingComma.ReplaceAll(cleaned, []byte("$1")), &entries); err != nil {
		return nil, fmt.Errorf("parse %s: %w", sourceURL, err)
	}

	byKey := make(map[string]string, len(entries)*2)
	for _, e := range entries {
		if e.URL == "" {
			continue
		}
		tpl := strings.ReplaceAll(e.URL, "{searchTerms}", "%s")
		for _, t := range e.Ts {
			if t == "" {
				continue
			}
			byKey[strings.ToLower(t)] = tpl
		}
	}
	return byKey, nil
}
