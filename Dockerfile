FROM node:24-alpine AS frontend
ENV NODE_OPTIONS=--use-openssl-ca
WORKDIR /app
RUN npm install -g pnpm@12.14.0
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY frontend/package.json ./frontend/
RUN pnpm install --frozen-lockfile --ignore-scripts
COPY frontend/ ./frontend/
RUN pnpm --filter frontend build

FROM golang:1.27.2-alpine AS builder
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
COPY --from=frontend /app/frontend/build /app/frontend/build
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/clasp ./cmd/clasp

FROM alpine:3.24
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/clasp .
COPY --from=builder /app/frontend/build ./frontend/build
EXPOSE 8080
CMD ["./clasp"]
