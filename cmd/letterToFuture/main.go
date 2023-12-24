package main

import (
	"flag"
	"log"

	commander "github.com/Destinyxus/botLetterToFuture/internal/bot_commander"
	"github.com/Destinyxus/botLetterToFuture/internal/config"
)

func main() {
	var path = flag.String("cfg-path", "/home/vladimir/GolandProjects/botLetterToFuture/internal/config/config.toml", "config path")

	flag.Parse()

	cfg, err := config.New(*path)
	if err != nil {
		log.Fatal(err)
	}

	botCommander := commander.New(commander.WithLogger(), commander.WithTgAPI(cfg.TelegramToken), commander.WithEmailSender(cfg.SendGridKey, cfg.LetterName, cfg.SendGridAddress))

	if err = botCommander.Start(); err != nil {
		log.Fatal(err)
	}
}
