# syntax=docker/dockerfile:1

FROM golang:1.22.1

WORKDIR /app

COPY go.mod go.sum ./

COPY ./ ./

RUN go mod download

RUN go build -o main ./cmd/

EXPOSE 8080

CMD ["/app/main"]