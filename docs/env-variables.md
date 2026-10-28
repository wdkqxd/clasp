# instance configuration (env variables)

if you have your own clasp instance, you can change its config through environment variables.

## general

| variable | default | example |
| --- | --- | --- |
| `PORT` | `8080` | `PORT=9000` |
| `STATIC_DIR` | `frontend/build` | `STATIC_DIR=/srv/clasp/frontend/build` |
| `CORS_ORIGIN` | `*` | `CORS_ORIGIN=https://search.example.com` |

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
