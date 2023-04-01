FROM golang

RUN go version

WORKDIR /bot

COPY . .

RUN go mod download

RUN go build -o telegramBot ./cmd/letterToFuture/main.go

EXPOSE 80

CMD ["./telegramBot"]

