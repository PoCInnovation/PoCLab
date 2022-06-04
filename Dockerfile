# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum /

RUN go mod download

COPY *.go ./

RUN go get Bot
RUN go build -o poc-lab

ENTRYPOINT ["./poc-lab"]