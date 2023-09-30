# Builder
FROM golang:1.19.4-alpine3.17 as builder

RUN apk update && apk upgrade && \
    apk --update add git make bash build-base

WORKDIR /app

# air installation
RUN go install github.com/cosmtrek/air@latest

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]

COPY . .

RUN make build

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app 

WORKDIR /app 

EXPOSE 9090

COPY --from=builder /app/engine /app/

CMD /app/engine