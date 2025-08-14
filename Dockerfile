FROM golang:1.24-bookworm@sha256:e617461712dbebf8768e10c1a5deab2833d67d2b692894cb8f4f0a3c19a8efb5 AS builder
ENV CGO_ENABLED=0
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY ./cmd/ ./cmd/

RUN go build -o restic-repo-exporter ./cmd/restic-repo-exporter

FROM debian:bookworm-slim@sha256:b1a741487078b369e78119849663d7f1a5341ef2768798f7b7406c4240f86aef
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
