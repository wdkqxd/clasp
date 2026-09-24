package wikipedia

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"metasearch/internal/engine"
)

type Article struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Extract     string `json:"extract,omitempty"`
	URL         string `json:"url"`
	Thumbnail   string `json:"thumbnail,omitempty"`
	Facts       []Fact `json:"facts,omitempty"`
}

type Fact struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type Client struct {
	http engine.HTTPClient
}

func NewClient(httpClient engine.HTTPClient) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 5 * time.Second}
	}
	return &Client{http: httpClient}
}

type apiResponse struct {
	Error *struct {
		Code string `json:"code"`
		Info string `json:"info"`
	} `json:"error"`
	Query struct {
		Pages map[string]struct {
			Title     string `json:"title"`
			Extract   string `json:"extract"`
			FullURL   string `json:"fullurl"`
			Thumbnail *struct {
				Source string `json:"source"`
			} `json:"thumbnail"`
			Terms struct {
				Description []string `json:"description"`
			} `json:"terms"`
			PageProps map[string]string `json:"pageprops"`
			Revisions []struct {
				Slots struct {
					Main map[string]string `json:"main"`
				} `json:"slots"`
			} `json:"revisions"`
		} `json:"pages"`
	} `json:"query"`
}

func (c *Client) Search(ctx context.Context, query, language string) (*Article, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, nil
	}

	var lastErr error
	for _, candidate := range languageCandidates(query, language) {
		article, err := c.searchLanguage(ctx, query, candidate)
		if err != nil {
			lastErr = err
			continue
		}
		if article != nil {
			return article, nil
		}
	}
	return nil, lastErr
}

func (c *Client) searchLanguage(ctx context.Context, query, language string) (*Article, error) {

	endpoint := fmt.Sprintf("https://%s.wikipedia.org/w/api.php", language)
	params := url.Values{
		"action":       {"query"},
		"format":       {"json"},
		"origin":       {"*"},
		"generator":    {"search"},
		"gsrsearch":    {query},
		"gsrnamespace": {"0"},
		"gsrlimit":     {"1"},
		"redirects":    {"1"},
		"prop":         {"extracts|info|pageimages|pageterms|revisions"},
		"exintro":      {"1"},
		"explaintext":  {"1"},
		"inprop":       {"url"},
		"pithumbsize":  {"320"},
		"wbptterms":    {"description"},
		"rvprop":       {"content"},
		"rvslots":      {"main"},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "metasearch/1.0 (Wikipedia API client)")

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wikipedia: unexpected status %s", res.Status)
	}

	var payload apiResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if payload.Error != nil {
		return nil, fmt.Errorf("wikipedia: %s", payload.Error.Info)
	}
	for _, page := range payload.Query.Pages {
		if page.Title == "" {
			continue
		}
		if _, disambiguation := page.PageProps["disambiguation"]; disambiguation {
			return nil, nil
		}
		if !isLikelyMatch(query, page.Title) {
			return nil, nil
		}
		article := &Article{
			Title:   page.Title,
			Extract: strings.TrimSpace(page.Extract),
			URL:     page.FullURL,
		}
		if len(page.Revisions) > 0 {
			content := page.Revisions[0].Slots.Main["*"]
			if content == "" {
				content = page.Revisions[0].Slots.Main["content"]
			}
			article.Facts = parseFacts(content)
		}
		if page.Thumbnail != nil {
			article.Thumbnail = page.Thumbnail.Source
		}
		if len(page.Terms.Description) > 0 {
			article.Description = strings.TrimSpace(page.Terms.Description[0])
		}
		if article.URL == "" {
			article.URL = "https://" + language + ".wikipedia.org/wiki/" + url.PathEscape(strings.ReplaceAll(page.Title, " ", "_"))
		}
		return article, nil
	}
	return nil, nil
}

func parseFacts(wikitext string) []Fact {
	start := strings.Index(strings.ToLower(wikitext), "{{infobox")
	if start < 0 {
		return nil
	}
	depth := 0
	end := len(wikitext)
	for i := start; i+1 < len(wikitext); i++ {
		switch wikitext[i : i+2] {
		case "{{":
			depth++
			i++
		case "}}":
			depth--
			i++
			if depth == 0 {
				end = i + 1
				i = len(wikitext)
			}
		}
	}

	lines := strings.Split(wikitext[start:end], "\n")
	facts := make([]Fact, 0, 6)
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "|") {
			continue
		}
		parts := strings.SplitN(strings.TrimSpace(strings.TrimPrefix(line, "|")), "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := cleanWikiValue(parts[1])
		if key == "" || value == "" || isTechnicalFact(key) {
			continue
		}
		facts = append(facts, Fact{Label: factLabel(key), Value: value})
		if len(facts) == 6 {
			break
		}
	}
	return facts
}

func cleanWikiValue(value string) string {
	value = regexp.MustCompile(`<!--.*?-->`).ReplaceAllString(value, "")
	value = regexp.MustCompile(`\[\[([^\]|]+)\|([^\]]+)\]\]`).ReplaceAllString(value, "$2")
	value = expandUbl(value)
	value = regexp.MustCompile(`\{\{[^{}]*\}\}`).ReplaceAllString(value, "")
	value = regexp.MustCompile(`\[\[([^\]]+)\]\]`).ReplaceAllString(value, "$1")
	value = regexp.MustCompile(`(?i)<ref[^>]*>.*?</ref>`).ReplaceAllString(value, "")
	value = regexp.MustCompile(`(?i)<ref[^/]*/>`).ReplaceAllString(value, "")
	value = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(value, "")
	return strings.TrimSpace(strings.Join(strings.Fields(value), " "))
}

func expandUbl(s string) string {
	needle := "{{ubl|"
	lower := strings.ToLower(s)
	for {
		idx := strings.Index(lower, needle)
		if idx < 0 {
			break
		}
		start := idx
		depth := 0
		end := start
		for i := start; i+1 < len(s); i++ {
			if s[i] == '{' && s[i+1] == '{' {
				depth++
				i++
			} else if s[i] == '}' && s[i+1] == '}' {
				depth--
				i++
				if depth == 0 {
					end = i + 1
					break
				}
			}
		}
		if end <= start {
			break
		}
		inner := s[start+len("{{") : end-len("}}")]
		if !strings.HasPrefix(strings.ToLower(inner), "ubl|") {
			lower = lower[:start+len("{{")] + lower[end:]
			s = s[:start+len("{{")] + s[end:]
			continue
		}
		inner = inner[len("ubl|"):]
		items := splitTopLevel(inner, '|')
		for i, item := range items {
			item = strings.TrimSpace(item)
			item = regexp.MustCompile(`\s*\(acting\)`).ReplaceAllString(item, "")
			items[i] = item
		}
		replacement := strings.Join(items, " / ")
		s = s[:start] + replacement + s[end:]
		lower = strings.ToLower(s)
	}
	return s
}

func splitTopLevel(s string, sep rune) []string {
	var parts []string
	braceDepth := 0
	linkDepth := 0
	cur := strings.Builder{}
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == '{' && i+1 < len(runes) && runes[i+1] == '{' {
			braceDepth++
		} else if r == '}' && i+1 < len(runes) && runes[i+1] == '}' {
			braceDepth--
		} else if r == '[' && i+1 < len(runes) && runes[i+1] == '[' {
			linkDepth++
		} else if r == ']' && i+1 < len(runes) && runes[i+1] == ']' {
			linkDepth--
		}
		if r == sep && braceDepth == 0 && linkDepth == 0 {
			parts = append(parts, cur.String())
			cur.Reset()
			continue
		}
		cur.WriteRune(r)
	}
	if cur.Len() > 0 {
		parts = append(parts, cur.String())
	}
	return parts
}

func isTechnicalFact(key string) bool {
	key = strings.ToLower(key)
	return strings.Contains(key, "image") || strings.Contains(key, "logo") || strings.Contains(key, "map") || strings.Contains(key, "caption") || strings.Contains(key, "module")
}

func factLabel(key string) string {
	key = strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(key), "_", " "), "-", " ")
	return cases.Title(language.English).String(strings.ToLower(key))
}

var languagePattern = regexp.MustCompile(`^[a-z]{2,3}$`)
var wordsPattern = regexp.MustCompile(`[\p{L}\p{N}]+`)

var ignoredQueryWords = map[string]struct{}{
	"a": {}, "an": {}, "and": {}, "define": {}, "is": {}, "of": {}, "the": {}, "what": {}, "who": {},
	"что": {}, "кто": {}, "такой": {}, "такое": {}, "определение": {}, "и": {}, "в": {}, "из": {}, "это": {},
	"что такое": {}, "как": {}, "quien": {}, "quién": {}, "qué": {}, "que": {}, "es": {}, "una": {}, "un": {}, "el": {}, "la": {}, "los": {}, "las": {},
	"was": {}, "were": {}, "how": {}, "wer": {}, "wie": {}, "ist": {}, "der": {}, "die": {}, "das": {},
}

func languageCandidates(query, requested string) []string {
	requested = strings.ToLower(strings.TrimSpace(requested))
	if requested != "" && requested != "auto" && languagePattern.MatchString(requested) {
		return []string{requested}
	}

	for _, r := range strings.ToLower(query) {
		switch {
		case r >= '\u0400' && r <= '\u04FF': // cyrillic
			if strings.ContainsAny(query, "іїєґ") {
				return []string{"uk", "ru", "en"}
			}
			if strings.ContainsAny(query, "ёы") {
				return []string{"ru", "uk", "en"}
			}
			return []string{"bg", "ru", "uk", "sr", "en"}
		case r >= '\u4E00' && r <= '\u9FFF', r >= '\u3400' && r <= '\u4DBF': // CJK Unified Ideographs
			return []string{"zh", "en"}
		case r >= '\u3040' && r <= '\u309F', r >= '\u30A0' && r <= '\u30FF': // Hiragana + Katakana
			return []string{"ja", "en"}
		case r >= '\uAC00' && r <= '\uD7AF': // hangul syllables
			return []string{"ko", "en"}
		case r >= '\u0600' && r <= '\u06FF': // arabic
			return []string{"ar", "en"}
		case r >= '\u0900' && r <= '\u097F': // devanagari
			return []string{"hi", "en"}
		case r >= '\u0E00' && r <= '\u0E5F': // thai
			return []string{"th", "en"}
		}
	}

	return []string{"en", "es", "fr", "de", "it", "pt", "nl", "pl", "tr"}
}

func isLikelyMatch(query, title string) bool {
	queryWords := wordsPattern.FindAllString(strings.ToLower(query), -1)
	titleWords := wordsPattern.FindAllString(strings.ToLower(title), -1)
	if len(queryWords) == 0 || len(titleWords) == 0 {
		return false
	}
	titleSet := make(map[string]struct{}, len(titleWords))
	for _, word := range titleWords {
		titleSet[word] = struct{}{}
	}
	found := 0
	for _, word := range queryWords {
		if _, ignored := ignoredQueryWords[word]; ignored {
			continue
		}
		if _, ok := titleSet[word]; !ok {
			return false
		}
		found++
	}
	return found > 0
}
