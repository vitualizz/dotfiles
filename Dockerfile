FROM golang:1.24-bookworm AS builder

WORKDIR /build

RUN apt-get update && apt-get install -y --no-install-recommends git \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o dotfiles ./cmd/dotfiles/

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    bash curl wget git zsh sudo \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /build/dotfiles /app/dotfiles

RUN chmod +x /app/dotfiles

CMD ["/app/dotfiles"]