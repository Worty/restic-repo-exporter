FROM golang:1.25-trixie@sha256:5e43883cbcb374df63de6fb695a86218a852b829055d2aa88c260a4189be46c5 AS builder
ENV CGO_ENABLED=0
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY ./cmd/ ./cmd/

RUN go build -o restic-repo-exporter ./cmd/restic-repo-exporter

FROM debian:trixie-slim@sha256:a347fd7510ee31a84387619a492ad6c8eb0af2f2682b916ff3e643eb076f925a
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates bzip2 curl && \
    rm -rf /var/lib/apt/lists/*

RUN curl -sSL https://github.com/restic/restic/releases/download/v0.18.0/restic_0.18.0_linux_amd64.bz2 | bunzip2 > /usr/local/bin/restic && \
    chmod +x /usr/local/bin/restic && \
    restic version && \
    restic self-update && \
    restic version

COPY --from=builder /app/restic-repo-exporter /app/restic-repo-exporter
ENTRYPOINT ["/app/restic-repo-exporter"]
