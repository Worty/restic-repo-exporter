FROM golang:1.25-bookworm@sha256:08c8ac4f6b588d7679ee764023f6ccbf12f629fe7f35d899eca3fcd71813e41e AS builder
ENV CGO_ENABLED=0
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY ./cmd/ ./cmd/

RUN go build -o restic-repo-exporter ./cmd/restic-repo-exporter

FROM debian:bookworm-slim@sha256:df52e55e3361a81ac1bead266f3373ee55d29aa50cf0975d440c2be3483d8ed3
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
