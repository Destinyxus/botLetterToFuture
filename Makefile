.PHONY: 
include .env

build:
	go build -o telegramBot cmd/letterToFuture/main.go

run: build
	./telegramBot

mock:
	mockgen -source=internal/bot_commander/bot-commander.go \
		-destination=internal/bot_commander/mocks/mocks.go 


