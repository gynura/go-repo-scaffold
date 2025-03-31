FROM golang:1.24.1 AS builder

WORKDIR /app

COPY app/go.mod app/go.sum ./
RUN go mod download

COPY app/ .

RUN go build -o bin/usersvc ./cmd/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/bin/usersvc /app/usersvc

COPY test/examples/service-config.yaml /app/config/service-config.yaml

EXPOSE 8080

CMD ["/app/usersvc", "--config-path", "/app/config"]
