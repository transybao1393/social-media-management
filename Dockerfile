FROM golang:1.19.5-alpine3.17
RUN apk add build-base vim
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

RUN go install github.com/cosmtrek/air@latest

COPY . ./