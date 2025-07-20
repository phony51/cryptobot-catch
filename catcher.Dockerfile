FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=linux CGO_ENABLED=0 GOAMD64=v3 go build  \
    -trimpath -ldflags="-s -w" \
    -gcflags="-d=ssa/check_bce=0 -m=0" \
    -o catcher.exe  \
    "./cmd/catch/main.go"

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/configuration.json /app
COPY --from=builder /app/sessions/ /app/sessions/
COPY --from=builder /app/bin/ /app/bin/

ENTRYPOINT ["/app/bin/catcher"]