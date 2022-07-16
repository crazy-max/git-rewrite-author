# syntax=docker/dockerfile:1

ARG GO_VERSION="1.18"
ARG GORELEASER_XX_VERSION="1.2.5"
ARG ALPINE_VERSION="3.16"

FROM --platform=$BUILDPLATFORM crazymax/goreleaser-xx:${GORELEASER_XX_VERSION} AS goreleaser-xx
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS base
ENV CGO_ENABLED=0
COPY --from=goreleaser-xx / /
RUN apk add --no-cache file git
WORKDIR /src

FROM base AS vendored
RUN --mount=type=bind,source=.,target=/src,rw \
  --mount=type=cache,target=/go/pkg/mod \
  go mod tidy && go mod download

FROM vendored AS build
ARG TARGETPLATFORM
RUN --mount=type=bind,target=. \
  --mount=type=cache,target=/root/.cache \
  --mount=target=/go/pkg/mod,type=cache \
  goreleaser-xx --debug \
    --name "git-rewrite-author" \
    --dist "/out" \
    --main="./cmd" \
    --flags="-trimpath" \
    --ldflags="-s -w -X 'main.version={{.Version}}'" \
    --files="CHANGELOG.md" \
    --files="LICENSE" \
    --files="README.md"

FROM scratch AS artifact
COPY --from=build /out/*.tar.gz /
COPY --from=build /out/*.zip /

FROM scratch AS binary
COPY --from=build /usr/local/bin/git-rewrite-author* /
