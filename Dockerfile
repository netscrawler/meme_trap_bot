FROM golang:1.26-alpine AS builder

RUN apk add --no-cache \
    gcc \
    musl-dev \
    git \
    sqlite-dev

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o bot ./cmd/bot

FROM alpine:3.23

RUN apk add --no-cache \
    ca-certificates \
    sqlite-libs \
    sqlite

WORKDIR /app

COPY --from=builder /build/bot .

RUN chmod 755 /app

ENTRYPOINT ["./bot"]
