package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	commander "github.com/Destinyxus/botLetterToFuture/internal/bot_commander"
	"github.com/Destinyxus/botLetterToFuture/internal/config"
	"github.com/Destinyxus/botLetterToFuture/internal/storage"
	"github.com/Destinyxus/botLetterToFuture/pkg/postgresconn"
)

func main() {
	var path = flag.String("cfg-path", "/home/vladimir/GolandProjects/botLetterToFuture/internal/config/config.toml", "config path")

	flag.Parse()

	cfg, err := config.New(*path)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	conn, err := postgresconn.New(ctx, "")
	if err != nil {
		log.Fatal(err)
	}

	st := storage.New(conn)

	botCommander := commander.New(st, commander.WithLogger(), commander.WithTgAPI(cfg.TelegramToken), commander.WithEmailSender(cfg.SendGridKey, cfg.LetterName, cfg.SendGridAddress))

	if err = botCommander.Start(ctx); err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
}
