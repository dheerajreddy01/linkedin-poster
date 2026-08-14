FROM golang:1.25-bookworm AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=1 GOOS=linux go build -o /out/server ./cmd/server

FROM debian:bookworm-slim
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=build /out/server ./server
COPY frontend ./frontend

ENV PORT=8081
ENV DB_PATH=/data/poster.db

EXPOSE 8081
CMD ["./server"]
