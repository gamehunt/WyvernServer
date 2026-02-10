FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main ./cmd/server

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
