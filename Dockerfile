FROM golang:1.27-trixie@sha256:0982f930de50a4f1a2b4453d51651f0031082ef2e3a25deb3c763fc39a1094a0 AS builder
ENV CGO_ENABLED=0
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY ./cmd/ ./cmd/

RUN go build -o restic-repo-exporter ./cmd/restic-repo-exporter

FROM debian:trixie-slim@sha256:a29215f6a35e51e22adffa17f89e9d2ef06214e64a2bad10d765c46aea49f11f
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
