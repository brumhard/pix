# syntax = docker/dockerfile:1.4

FROM --platform=$BUILDPLATFORM golang:1.20 AS base-go
ARG TARGETOS TARGETARCH
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
WORKDIR /src
COPY go.mod go.sum .
RUN --mount=type=cache,target=/go/pkg \
    go mod download

FROM --platform=$BUILDPLATFORM base-go AS build-go
COPY pkg/ pkg/
# COPY cmd/ cmd/
COPY main.go .
RUN --mount=type=cache,target=/go/pkg \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o /app ./main.go

FROM gcr.io/distroless/static:nonroot
COPY --from=build-go /app /
# this is the numeric version of user nonroot:nonroot to check runAsNonRoot in kubernetes
USER 65532:65532
ENTRYPOINT ["/app"]
