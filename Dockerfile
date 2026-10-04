FROM denoland/deno:alpine AS frontend
WORKDIR /app
COPY frontend/deno.json frontend/deno.lock ./
RUN deno install --frozen
COPY frontend/ ./
RUN deno task build

FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/build ./frontend/build
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o clasp .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/clasp .
COPY --from=frontend /app/build ./frontend/build
EXPOSE 8080
CMD ["./clasp"]
