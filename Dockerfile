# 빌드 스테이지
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download && \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main ./cmd/main.go

# 실행 스테이지
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

RUN chmod +x ./main

EXPOSE ${PORT}

CMD ["./main"]
