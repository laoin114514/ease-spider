# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder
WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o spider .

FROM alpine:3.20
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app && adduser -S app -G app

COPY --from=builder /app/spider /app/spider

USER app
CMD ["./spider"]