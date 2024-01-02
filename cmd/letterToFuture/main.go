package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Destinyxus/botLetterToFuture/internal/storage"
	"github.com/Destinyxus/botLetterToFuture/pkg/postgresconn"
	"log"
	"os/signal"
	"sync"
	"syscall"

	commander "github.com/Destinyxus/botLetterToFuture/internal/bot_commander"
	"github.com/Destinyxus/botLetterToFuture/internal/config"
)

func main() {
	var path = flag.String("cfg-path", "internal/config/config.toml", "config path")

	flag.Parse()

	cfg, err := config.New(*path)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	conn, err := postgresconn.New(ctx, cfg.PostgresAddress)
	if err != nil {
		log.Fatal(err)
	}

	st, err := storage.New(conn)
	if err != nil {
		log.Fatal(err)
	}

	botCommander := commander.New(
		st,
		commander.WithLogger(), commander.WithTgAPI(cfg.TelegramToken), commander.WithEmailSender(cfg.SendGridKey, cfg.LetterName, cfg.SendGridAddress))

	var wg sync.WaitGroup

	if err = botCommander.Start(ctx, &wg); err != nil {
		log.Fatal(err)
	}

	wg.Wait()

	fmt.Println("graceful shutdown")
}
