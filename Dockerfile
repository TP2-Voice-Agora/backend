# syntax=docker/dockerfile:1
FROM golang:1.23.8-alpine

WORKDIR /agora

COPY go.mod go.sum ./

RUN go mod download

COPY . .    

RUN go build -o main ./app

EXPOSE 8080

CMD ["./main"]