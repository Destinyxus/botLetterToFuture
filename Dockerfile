FROM golang:1.20-alpine3.16 AS builder
RUN go version


COPY . /github.com/Destinyxus/botLetterToFuture/
WORKDIR /github.com/Destinyxus/botLetterToFuture/

RUN go mod download
RUN go build -o ./bin/bot ./cmd/letterToFuture/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/Destinyxus/botLetterToFuture/bin/bot .
COPY --from=0 /github.com/Destinyxus/botLetterToFuture/configs configs/

EXPOSE 80

CMD ["./bot"]