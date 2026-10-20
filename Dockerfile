FROM node:24-alpine AS frontend
WORKDIR /app
RUN npm install -g pnpm@12.11.2
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY frontend/package.json ./frontend/
RUN pnpm install --frozen-lockfile --ignore-scripts
COPY frontend/ ./frontend/
RUN pnpm --filter frontend build

FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/frontend/build ./frontend/build
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o clasp .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/clasp .
COPY --from=builder /app/frontend/build ./frontend/build
EXPOSE 8080
CMD ["./clasp"]
