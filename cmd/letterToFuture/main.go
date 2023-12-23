package main

import (
	"log"

	commander "github.com/Destinyxus/botLetterToFuture/internal/bot_commander"
)

func main() {
	var (
	//path = flag.String("cfg-path", "/home/vladimir/GolandProjects/botLetterToFuture/internal/config/config.toml", "")
	)

	botCommander := commander.New(commander.WithLogger(), commander.WithTgAPI())

	if err := botCommander.Start(); err != nil {
		log.Fatal(err)
	}
}
