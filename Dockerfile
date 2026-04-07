# ---- Build stage ----
FROM golang:1.25.3-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /build/app ./cmd

# ---- Run stage ----
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/app .
COPY config/staging.yaml config/staging.yaml
COPY config/production.yaml config/production.yaml

CMD ["./app", "run"]
