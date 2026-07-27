# -----------------------------
# Build Stage
# -----------------------------
FROM golang:1.26-alpine AS builder

ARG NAME
ARG PATH

WORKDIR /app
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build \
    -ldflags="-s -w" \
    -o ${NAME} \
    ${PATH}

# -----------------------------
# Runtime Stage
# -----------------------------
FROM alpine:3.22

ARG NAME
ARG PORT

WORKDIR /app
RUN apk add --no-cache ca-certificates

ARG EXE_PATH=/app/${NAME}
COPY --from=builder ${EXE_PATH} .

EXPOSE ${PORT}

CMD ["${EXE_PATH} -app=${NAME}"]
