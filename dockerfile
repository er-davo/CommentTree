FROM golang:alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /comment-tree

COPY app/go.mod /comment-tree/

RUN go mod download

COPY app/ /comment-tree/

RUN go build -o build/main cmd/main.go

FROM alpine:latest AS runner

WORKDIR /app

COPY --from=builder /comment-tree/build/main /app/

COPY /config/config.yaml /app/config.yaml
COPY public/ /app/public/
COPY /migrations /app/migrations

ENV CONFIG_PATH=/app/config.yaml
ENV APP_MIGRATION_DIR=/app/migrations

CMD [ "/app/main" ]