FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/main "./cmd/catch/main.go"

FROM alpine:3.19

WORKDIR /app

VOLUME ["/app/sessions", "/app/configuration.json"]

COPY --from=builder /app/bin /app/sessions /app/configuration.json /app/
ENTRYPOINT ["/app/bin/main"]