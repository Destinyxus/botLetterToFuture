.PHONY:
include .env

build:
	go build -o ./bin/bot cmd/letterToFuture/main.go

run: build
	./bin/bot

build-image:
		docker build -t futureletter:1.0 .

start-container:
	docker run  -p 80:80 letter:0.1
