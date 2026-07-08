FROM golang:1.26.4-bookworm AS deps

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

FROM golang:1.26.4-bookworm AS builder

WORKDIR /app

COPY --from=deps /go/pkg /go/pkg
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -ldflags="-w -s" -o shopify-storefront-api ./cmd/api

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/shopify-storefront-api .
COPY bin/start_server.sh bin/start_server.sh

RUN chmod +x /app/bin/start_server.sh

EXPOSE 8888

CMD ["bin/start_server.sh"]
