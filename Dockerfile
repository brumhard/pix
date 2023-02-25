# syntax = docker/dockerfile:1.4

FROM --platform=$BUILDPLATFORM ubuntu:22.04 AS base-flutter
# required by flutter sdk
RUN apt update && apt install -y curl git unzip xz-utils zip libglu1-mesa
# Set up new user
RUN useradd -ms /bin/bash developer
USER developer
WORKDIR /home/developer
# Download Flutter SDK
RUN git clone --depth=1 --branch=stable https://github.com/flutter/flutter.git
ENV PATH "$PATH:/home/developer/flutter/bin"
# Run basic check to download Dark SDK
RUN flutter doctor

FROM --platform=$BUILDPLATFORM base-flutter AS build-flutter
WORKDIR /frontend
# somehow flutter pub get tries to update pubspec.lock every time
COPY --chown=developer ./frontend .
RUN --mount=type=cache,target=/root/.pub-cache \
    flutter pub get && flutter build web

FROM --platform=$BUILDPLATFORM golang:1.20 AS base-go
WORKDIR /src
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ARG TARGETOS TARGETARCH
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
COPY go.mod go.sum .
RUN --mount=type=cache,target=/go/pkg \
    go mod download

FROM --platform=$BUILDPLATFORM base-go AS build-go
COPY --from=build-flutter /frontend/build/web frontend/build/web
COPY frontend/embed.go frontend/
COPY pkg/ pkg/
COPY main.go .
RUN --mount=type=cache,target=/go/pkg \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o /app ./main.go

FROM gcr.io/distroless/static:nonroot
COPY --from=build-go /app /
# this is the numeric version of user nonroot:nonroot to check runAsNonRoot in kubernetes
USER 65532:65532
ENTRYPOINT ["/app"]
