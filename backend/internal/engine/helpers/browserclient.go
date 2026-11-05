package helpers

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	fhttp "github.com/bogdanfinn/fhttp"
	fcookiejar "github.com/bogdanfinn/fhttp/cookiejar"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"

	"clasp/internal/engine"
)

type BrowserClient struct {
	inner tls_client.HttpClient
}

func NewBrowserClient() (*BrowserClient, error) {
	jar, err := fcookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("cookie jar: %w", err)
	}

	opts := []tls_client.HttpClientOption{
		tls_client.WithClientProfile(profiles.Chrome_150),
		tls_client.WithTimeoutSeconds(25),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar),
	}
	if proxy := strings.TrimSpace(os.Getenv("PROXY_URL")); proxy != "" {
		opts = append(opts, tls_client.WithProxyUrl(proxy))
	}

	c, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), opts...)
	if err != nil {
		return nil, fmt.Errorf("tls_client: %w", err)
	}
	return &BrowserClient{inner: c}, nil
}

func (b *BrowserClient) Do(req *http.Request) (*http.Response, error) {
	freq, err := fhttp.NewRequest(req.Method, req.URL.String(), req.Body)
	if err != nil {
		return nil, err
	}
	for k, vs := range req.Header {
		for _, v := range vs {
			freq.Header.Add(k, v)
		}
	}
	freq.Host = req.Host
	if freq.Host == "" {
		freq.Host = req.URL.Host
	}

	fresp, err := b.inner.Do(freq)
	if err != nil {
		return nil, err
	}

	return &http.Response{
		Status:     fresp.Status,
		StatusCode: fresp.StatusCode,
		Proto:      fresp.Proto,
		ProtoMajor: fresp.ProtoMajor,
		ProtoMinor: fresp.ProtoMinor,
		Header:     http.Header(fresp.Header),
		Body:       fresp.Body,
		Request:    req,
	}, nil
}

var _ engine.HTTPClient = (*BrowserClient)(nil)
