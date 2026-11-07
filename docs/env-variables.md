# instance configuration (env variables)

if you have your own clasp instance, you can change its config through environment variables.

## general

| variable | default | example |
| --- | --- | --- |
| `PORT` | `8080` | `PORT=9000` |
| `STATIC_DIR` | `frontend/build` | `STATIC_DIR=/srv/clasp/frontend/build` |
| `CORS_ORIGIN` | `*` | `CORS_ORIGIN=https://search.example.com` |
| `PROXY_URL` | *(empty)* | `PROXY_URL=http://user:pass@proxy:8080` |

## upstream protection

search engines rate limit by ip, so a single instance shares one budget between all its users.
these settings reduce how often clasp has to ask an engine in the first place.

| variable | default | example |
| --- | --- | --- |
| `SEARCH_CACHE_TTL` | `300` | `SEARCH_CACHE_TTL=600` |
| `SEARCH_CACHE_MAX_ITEMS` | `1024` | `SEARCH_CACHE_MAX_ITEMS=4096` |
| `ENGINE_MIN_INTERVAL` | `1200` | `ENGINE_MIN_INTERVAL=2500` |
| `ENGINE_COOLDOWN` | `600` | `ENGINE_COOLDOWN=1800` |
| `ENGINE_COOLDOWN_MAX` | `3600` | `ENGINE_COOLDOWN_MAX=7200` |

## rate limiter

| variable | default | example |
| --- | --- | --- |
| `LIMITER_ENABLED` | `true` | `LIMITER_ENABLED=false` |
| `LIMITER_IPV4_SUBNET` | `24` | `LIMITER_IPV4_SUBNET=32` |
| `LIMITER_IPV6_SUBNET` | `56` | `LIMITER_IPV6_SUBNET=64` |
| `LIMITER_SCALE` | `1` | `LIMITER_SCALE=0.5` |
| `LIMITER_TRUSTED_PROXIES` | `127.0.0.0/8,::1` | `LIMITER_TRUSTED_PROXIES=10.0.0.0/8` |
| `LIMITER_EXEMPT_PATHS` | `/api/healthz` | `LIMITER_EXEMPT_PATHS=/api/healthz,/api/bangs` |

# notes

## general

### port

the port the server listens on. it is the only thing worth setting on a home instance. the vite dev
server also reads it and proxies `/api/*` to this port, so `pnpm dev` needs no extra configuration.

### static_dir

the directory with the built frontend that is served by the go server. the default
`frontend/build` is what `pnpm build` produces and what the docker image contains, so you usually do
not need to change it.

### cors_origin

the value of `access-control-allow-origin`. it is not protection, it only decides who a browser
will let read api responses from another origin, so `*` is fine unless you call the api from
another site.

### environment file

clasp reads `.env` from the working directory on startup. the pnpm scripts run the backend from
the repository root, so put `.env` next to `package.json`. a malformed file is reported in the log,
and a missing file is ignored. values from `.env` never override variables that are already set in
the environment, so docker, systemd and shell configuration always win.

### proxy_url

an optional proxy for all upstream requests. useful when your host ip is already flagged by the
search engines. only one proxy is supported for now; a rotating residential gateway is the
practical option for a public instance.

## upstream protection

### search_cache_ttl

how long an aggregated result is cached, in seconds. repeated and concurrent identical queries are
served from memory instead of hitting the engines again. `0` disables the cache. this is the
cheapest way to cut upstream traffic, because one user reloading or paginating a result page no
longer counts as a new search.

### search_cache_max_items

maximum number of cached responses before the cache is pruned or reset. `0` means unlimited.

### engine_min_interval

minimum gap between two requests to the same engine, in milliseconds, with random jitter. this
keeps a burst of clicks from looking like a script. raise it if you still get blocked; lower it if
you want snappier searches on a quiet instance.

### engine_cooldown / engine_cooldown_max

when an engine answers with `403`, `429` or duckduckgo's `202` anomaly page, clasp stops calling it
for `ENGINE_COOLDOWN` seconds and doubles that time on every further block, up to
`ENGINE_COOLDOWN_MAX`. other engines keep serving results in the meantime. without this, retrying a
blocked engine immediately is what turns a short throttle into a long ban.

## rate limiter

### limiter_enabled

turns ip-based rate limiting on or off. requests to the api are limited by both a short burst
window and a longer sustained window, with separate budgets for individual endpoints.

### limiter_ipv4_subnet / limiter_ipv6_subnet

clients are grouped into subnets before counting requests, so rotating addresses inside the same
network does not bypass the limit. the values are clamped to `0`–`32` for ipv4 and `0`–`128` for
ipv6. `32` / `128` counts every address separately, while `0` puts all clients of that family into
one bucket.

### limiter_scale

multiplies every rate limit, both burst and sustained. `0.5` halves the allowed request counts,
`2` doubles them. values less than or equal to zero are ignored.

### limiter_trusted_proxies

a comma-separated list of ip addresses and cidr ranges whose forwarding headers are trusted. when
the direct peer is in this list, clasp reads the client address from `x-forwarded-for` or
`x-real-ip`. behind a reverse proxy, set this to the proxy network; otherwise every caller will
look like the proxy and share one limit. keep the default on an instance exposed directly to the
internet.

### limiter_exempt_paths

a comma-separated list of exact paths that bypass rate limiting. `options` requests are always
allowed and are not affected by this setting.
