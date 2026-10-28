# self-hosting

clasp runs as an oci container or from source. the image already contains the built frontend.

the examples use docker; podman and other oci runtimes work the same way.

## container

```sh
git clone https://github.com/wdkq/clasp.git
cd clasp
docker build -t clasp .
docker run -d --name clasp --restart unless-stopped -p 8080:8080 clasp
```

open `http://localhost:8080`. `PORT` defaults to `8080`; if you change it, publish the same port and
set it with `-e PORT=...` or `--env-file .env`.

## compose

```yaml
services:
  clasp:
    build: .
    restart: unless-stopped
    env_file:
      - .env
    ports:
      - "${PORT:-8080}:${PORT:-8080}"
```

```sh
docker compose up -d --build
```

`env_file` is optional; remove it if you configure the service with `environment` instead.

## from source

requires go 1.26, node 24 and pnpm 12.

```sh
pnpm install
pnpm build
go run ./backend/cmd/clasp
```

run from the repository root so `.env` and `frontend/build` resolve correctly. use `pnpm dev` for
the dev server.

## configuration

all variables are listed in [environment variables](env-variables.md). `.env` is optional, is read
from the working directory, and never overrides variables that are already set.

## reverse proxy

terminate tls at a reverse proxy. add its addresses to `LIMITER_TRUSTED_PROXIES`; otherwise every
client looks like the proxy and shares one rate limit.

```caddyfile
search.example.com {
    reverse_proxy 127.0.0.1:8080
}
```

## updates

```sh
git pull
docker compose up -d --build
```
