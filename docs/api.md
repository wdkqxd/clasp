# api

clasp exposes a read-only json api under `/api/`. all endpoints use `get`; `options` returns `204`
for cors. responses are json unless noted otherwise.

`access-control-allow-origin` comes from `CORS_ORIGIN` and defaults to `*`.

## endpoints

| endpoint | description |
| --- | --- |
| `GET /api/search` | web search |
| `GET /api/images` | image search |
| `GET /api/suggest` | autocomplete |
| `GET /api/wikipedia` | wikipedia summary |
| `GET /api/bangs` | available bangs |
| `GET /api/healthz` | health and uptime |

## search

```http
GET /api/search?q=rust&limit=10&offset=0
```

| parameter | default | description |
| --- | --- | --- |
| `q` | required | search query |
| `limit` | `10` | results per page, `1`–`50` |
| `offset` | `0` | zero-based page. `page` is an alias |
| `engines` | all | comma-separated: `brave`, `ddg` (alias `duckduckgo`) |
| `timeout` | fast mode | integer `5`–`10` seconds or `off` |

unknown engine names are ignored. without `timeout`, engines get 3 seconds to respond and the whole
request is capped at 10 seconds; `off` gives them all 10 seconds.

```json
{
  "query": "rust",
  "took_ms": 412,
  "results": [
    {
      "title": "Rust Programming Language",
      "url": "https://www.rust-lang.org/",
      "snippet": "A language empowering everyone to build reliable software.",
      "engine": "brave",
      "engines": ["brave", "ddg"]
    }
  ],
  "engines": {
    "brave": { "status": "ok", "count": 10, "took_ms": 187 },
    "ddg": { "status": "ok", "count": 10, "took_ms": 301 }
  },
  "has_more": true,
  "infobox": { "title": "Rust", "url": "https://en.wikipedia.org/wiki/Rust_(programming_language)" }
}
```

| field | description |
| --- | --- |
| `results` | ranked and deduplicated results: `title`, `url`, `snippet`, `engine`, optional `engines`, optional `sitelinks` |
| `engines` | per-engine status: `status`, `count`, `took_ms`, optional `error` |
| `has_more` | best-effort hint that another page may exist |
| `infobox` | optional `{ title, url }` |
| `discussions` | optional list of `{ title, url, snippet? }` |
| `rich` | optional brave answers, currently cryptocurrency and stock |

an engine failure does not fail the request; it is reported in `engines` with `status: "error"`.

## images

```http
GET /api/images?q=red+panda&limit=20&offset=0
```

parameters are the same as for search, with `limit` defaulting to `20` (max `200`) and engines
`brave`, `ddg` and `startpage`. results from different engines are interleaved.

```json
{
  "query": "red panda",
  "took_ms": 680,
  "results": [
    {
      "title": "Red panda",
      "url": "https://example.com/red-panda.jpg",
      "source_url": "https://example.com/red-panda",
      "thumbnail": "https://example.com/red-panda-thumb.jpg",
      "width": 1920,
      "height": 1280,
      "engine": "brave",
      "engines": ["brave", "startpage"]
    }
  ],
  "engines": {
    "brave": { "status": "ok", "count": 200, "took_ms": 420 }
  },
  "has_more": true
}
```

image results contain `title`, `url`, optional `source_url`, `thumbnail`, `width`, `height`,
`engine` and optional `engines`. `results` is `null` when nothing is found.

## suggest

```http
GET /api/suggest?q=rust&engine=kagi
```

| parameter | default | description |
| --- | --- | --- |
| `q` | required | partial query |
| `engine` | `kagi` | `kagi`, `ddg` (alias `duckduckgo`), `brave` or `startpage` |

```json
{
  "query": "rust",
  "engine": "kagi",
  "suggestions": ["rust", "rust programming", "rust vs go"]
}
```

returns up to 7 suggestions. results are cached for 5 minutes. an upstream failure returns `502`
with `{ "engine": "...", "error": "..." }`; an unsupported engine returns `400`.

## wikipedia

```http
GET /api/wikipedia?q=rust&lang=en
```

| parameter | default | description |
| --- | --- | --- |
| `q` | required | article search query |
| `lang` | `auto` | two- or three-letter language code, or `auto` |

```json
{
  "title": "Rust (programming language)",
  "description": "general-purpose programming language",
  "extract": "Rust is a general-purpose programming language...",
  "url": "https://en.wikipedia.org/wiki/Rust_(programming_language)",
  "thumbnail": "https://upload.wikimedia.org/..."
}
```

`description`, `extract` and `thumbnail` are optional. when no article is found the endpoint
returns `200` with `null`. an upstream failure returns `502` with
`{ "error": "wikipedia unavailable" }`.

## bangs

```http
GET /api/bangs
```

returns a map of lowercase triggers to url templates. `%s` is replaced by the query on the client.

```json
{
  "g": "https://www.google.com/search?q=%s",
  "yt": "https://www.youtube.com/results?search_query=%s"
}
```

the list is loaded at startup and kept for the lifetime of the process.

## health

```http
GET /api/healthz
```

```json
{
  "status": "ok",
  "uptime": "1h2m3s"
}
```

the endpoint is exempt from rate limiting by default.

## errors

most errors are json with an `error` field. `429` is plain text.

| status | meaning |
| --- | --- |
| `400` | missing or invalid parameters |
| `405` | method other than `get` or `options` |
| `429` | rate limit exceeded |
| `502` | upstream autocomplete or wikipedia failure |

a `429` response includes `retry-after` and `ratelimit-reset` headers.

## rate limits

limits are applied per client subnet and can be changed with the `LIMITER_*` variables in
[environment variables](env-variables.md). defaults:

| endpoint | burst | sustained |
| --- | --- | --- |
| `/api/search` | 10 / 10s | 40 / 1m |
| `/api/images` | 10 / 10s | 60 / 1m |
| `/api/wikipedia` | 10 / 10s | 60 / 1m |
| `/api/suggest` | 40 / 2s | 300 / 1m |
| `/api/bangs` | 30 / 2s | 120 / 1m |
| other `/api/*` | 20 / 10s | 120 / 1m |
