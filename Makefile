.PHONY:
build:
	go build -o telegramBot cmd/letterToFuture/main.go

.PHONY:
run: build
	./telegramBot

.PHONY:
mock:
	mockgen -source=internal/bot_commander/bot-commander.go \
		-destination=internal/bot_commander/mocks/mocks.go 


