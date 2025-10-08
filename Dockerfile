FROM golang:1.24.5-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o ./bin/auth_server cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root

COPY --from=builder app/bin/auth_server .

CMD ["./auth_server"]