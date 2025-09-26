# syntax=docker/dockerfile:1

FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o pokerbot ./cmd/bot

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app

COPY --from=builder /app/pokerbot ./pokerbot

ENV TELEGRAM_BOT_TOKEN=""

ENTRYPOINT ["/app/pokerbot"]
