
# Build stage
FROM golang:1.27.4-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /geotrack \
    ./cmd

# Runtime stage
FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /geotrack /usr/local/bin/geotrack

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/geotrack"]
 