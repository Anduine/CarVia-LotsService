FROM golang:1.24-alpine3.21
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY configs/container.config.yaml configs/config.yaml
RUN go build -o /cmd/app/main ./cmd/app/main.go
CMD [ "/cmd/app/main" ]