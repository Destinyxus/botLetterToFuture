FROM golang:1.21.5-alpine3.19

ENV TZ=Europe/Moscow

RUN apk add tzdata
RUN cp /usr/share/zoneinfo/Europe/Moscow  /etc/localtime
RUN echo "Europe/Moscow" >  /etc/timezone

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build ./cmd/letterToFuture

CMD ["/app/letterToFuture"]