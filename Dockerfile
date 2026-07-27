# -----------------------------
# Build Stage
# -----------------------------
FROM golang:1.26-alpine AS builder

WORKDIR /src
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/ecom ./cmd

# -----------------------------
# Runtime Stage
# -----------------------------
FROM alpine:3.22

RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /out/ecom /app/ecom
# Optional: app env files for LoadForApp when not injected by Compose.
# Note: .dockerignore may exclude .env paths; Compose env_file injects process env instead.
COPY build/ /app/build/

EXPOSE 8080
ENTRYPOINT ["/app/ecom"]
CMD ["-app", "gateway"]
