# syntax=docker/dockerfile:1.7

FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder
WORKDIR /spider

RUN apk add --no-cache ca-certificates

ENV GOPROXY=https://goproxy.cn,direct \
    GOSUMDB=sum.golang.google.cn

# go.mod 用 replace 把内化的 Codeforces 客户端重定向到 pkg/codeforcesAPIClient
# （git submodule）。go mod download 需要先读到它的 go.mod，所以单独拷一份；
# 源码目录稍后整块 COPY。子模块没检出时这里会直接失败，属于预期（见 README 运行方式）。
COPY go.mod go.sum ./
COPY pkg/codeforcesAPIClient/go.mod pkg/codeforcesAPIClient/go.mod
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} go build -trimpath -ldflags="-s -w" -o spider ./cmd

FROM alpine:3.20
WORKDIR /spider

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app && adduser -S app -G app

COPY --from=builder /spider /spider
RUN mkdir -p /spider/data /spider/logs /config \
    && chown -R app:app /spider /config

VOLUME ["/spider/data", "/spider/logs", "/config"]

USER app
CMD ["./spider"]