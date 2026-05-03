FROM golang:1.24-alpine3.23 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/ ./cmd/
COPY ./pkg/ ./pkg/
COPY ./docs/ ./docs/
COPY ./internal/ ./internal/
COPY ./configs/container.config.yaml ./configs/config.yaml

RUN CGO_ENABLED=0 GOOS=linux go build -o lots_bin ./cmd/app/main.go

FROM alpine:3.23.3

WORKDIR /app

COPY --from=builder /app/lots_bin .
COPY --from=builder /app/configs/config.yaml /app/configs/config.yaml

CMD [ "./lots_bin" ]