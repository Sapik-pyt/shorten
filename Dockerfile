FROM golang:1.20 AS builder
WORKDIR /build

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o app ./cmd/main.go

CMD ["/build/./app"]