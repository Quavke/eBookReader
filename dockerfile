FROM golang:1.24.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# СТАТИЧЕСКАЯ сборка
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main

# Используем статический distroless
FROM gcr.io/distroless/static

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/cmd/main/config.yaml ./config.yaml

EXPOSE 8080

CMD ["./main"]