FROM golang:1.25.3-alpine

WORKDIR /app

RUN apk add git
RUN go install github.com/air-verse/air@v1.63.0

COPY go.mod go.sum ./
RUN go mod download

COPY .air.toml ./

CMD ["air", "-c", ".air.toml"]