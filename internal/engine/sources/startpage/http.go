package startpage

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"metasearch/internal/engine"
)

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
		return "", errors.New(resp.Status)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, engine.MaxBodyReadLimit))
	if err != nil {
		return "", err
	}
	return string(body), nil
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
