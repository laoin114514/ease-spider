# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder
WORKDIR /spider

RUN apk add --no-cache ca-certificates

ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o spider .

FROM alpine:3.20
WORKDIR /spider

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app && adduser -S app -G app

COPY --from=builder /spider/spider /spider/spider
COPY --from=builder /spider/config /spider/config

USER app
CMD ["./spider"]