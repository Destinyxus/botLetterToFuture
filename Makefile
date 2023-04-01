.PHONY:
include .env

build:
	go build -o telegramBot cmd/letterToFuture/main.go

run: build
	./telegramBot

