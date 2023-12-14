# Multi-stage build
FROM golang:1.17 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/task-queue cmd/main.go

# Minimal runtime image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/task-queue /app/task-queue

CMD ["./task-queue"]
