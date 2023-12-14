# İlk aşama: Bağımlılıkları çöz
FROM golang:1.21.5-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/task-queue cmd/main.go

# İkinci aşama: Minimal runtime imajı oluştur
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/task-queue /app/task-queue

EXPOSE 8080

CMD ["/app/task-queue"]
