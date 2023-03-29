build:
	go build -v ./cmd/letterToFuture
.DEFAULT_GOAL := build

run:build
	@./letterToFuture.exe