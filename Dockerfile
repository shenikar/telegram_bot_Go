# Builder stage
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go build -o telegram_bot cmd/main.go

# Final stage
FROM ubuntu:22.04
WORKDIR /root
COPY --from=builder /app/telegram_bot /root/telegram_bot
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/migrations /root/migrations
RUN chmod +x /root/telegram_bot
RUN apt-get update && apt-get install -y ca-certificates
CMD ["/root/telegram_bot"]