FROM golang:1.24.5-alpine AS builder

COPY . /github.com/sborsh1kmusora/auth/source
WORKDIR /github.com/sborsh1kmusora/auth/source

RUN go mod download
RUN go build -o ./bin/auth_server cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/sborsh1kmusora/auth/source/bin/auth_server .

CMD ["./auth_server"]