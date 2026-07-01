FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build \
    -ldflags="-s -w" \
    -o document-processor \
    ./cmd/main.go

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates

RUN addgroup -S app && adduser -S app -G app

COPY --from=builder /app/document-processor .
COPY --from=builder /app/config ./config

USER app

EXPOSE 8081
EXPOSE 9091

ENTRYPOINT ["./document-processor"]