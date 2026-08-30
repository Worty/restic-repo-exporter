FROM golang:1.27-trixie@sha256:df98008ecd2b0ecc9f0a94d1b07e3564a9c92b555369b33d9b5f60d0765b2db7 AS builder
ENV CGO_ENABLED=0
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY ./cmd/ ./cmd/

RUN go build -o restic-repo-exporter ./cmd/restic-repo-exporter

FROM debian:trixie-slim@sha256:d7e12182ce18b85b93007c1dedf31f2d29e01ccf3182cc4017c709b6259bc132
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates bzip2 curl && \
    rm -rf /var/lib/apt/lists/*

RUN curl -sSL https://github.com/restic/restic/releases/download/v0.19.0/restic_0.19.0_linux_amd64.bz2 | bunzip2 > /usr/local/bin/restic && \
    chmod +x /usr/local/bin/restic && \
    restic version && \
    restic self-update && \
    restic version

COPY --from=builder /app/restic-repo-exporter /app/restic-repo-exporter
ENTRYPOINT ["/app/restic-repo-exporter"]
