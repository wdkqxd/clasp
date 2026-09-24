package startpage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var reAnubisChallenge = regexp.MustCompile(`(?s)<script id="anubis_challenge" type="application/json">\s*(\{.*?\})\s*</script>`)

func hasChallenge(doc string) bool {
	return reAnubisChallenge.MatchString(doc)
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