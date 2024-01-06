package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/Destinyxus/botLetterToFuture/internal/storage"
	"github.com/Destinyxus/botLetterToFuture/pkg/postgresconn"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

	conn, err := postgresconn.New(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	st, err := storage.New(conn)
	if err != nil {
		log.Fatal(err)
	}

	botCommander, err := commander.New(
		st,
		*cfg,
		commander.WithLogger(),
		commander.WithTgAPI(cfg.TelegramToken),
		commander.WithEmailSender(cfg.EmailSender.EmailToken, cfg.EmailSender.ClientEmail, cfg.EmailSender.HostEmail, cfg.EmailSender.SMTPAddress),
	)
	if errors.Is(err, commander.DateIndexesError) {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	if err = botCommander.Start(ctx, &wg); err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(time.Minute)

loop:
	for {
		select {
		case <-ticker.C:
			if err = botCommander.CheckForActualDate(); err != nil {
				log.Fatal(err)
			}

			ticker.Reset(time.Minute)

		case <-ctx.Done():
			fmt.Println("graceful shutdown")

			break loop
		}
	}

	wg.Wait()
}
