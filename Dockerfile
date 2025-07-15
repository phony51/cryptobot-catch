FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
-ldflags="-s -w" \
-trimpath \
-o /app/bin/main \
./cmd/catch/main.go

FROM alpine:3.19

WORKDIR /app

VOLUME ["/app/sessions", "/app/configuration.json", "logs.json"]

COPY --from=builder /app/bin /app/sessions /app/configuration.json ./
ENTRYPOINT ["/app/bin/main"]