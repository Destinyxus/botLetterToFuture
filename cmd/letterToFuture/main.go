package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	commander "github.com/Destinyxus/botLetterToFuture/internal/bot_commander"
	"github.com/Destinyxus/botLetterToFuture/internal/config"
	"github.com/Destinyxus/botLetterToFuture/internal/emailSender"
	"github.com/Destinyxus/botLetterToFuture/internal/logger"
	"github.com/Destinyxus/botLetterToFuture/internal/storage"
	"github.com/Destinyxus/botLetterToFuture/pkg/postgresconn"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
}

func main() {
	path := *flag.String("cfg-path", "internal/config/config.toml", "config path")
	flag.Parse()

	cfg, err := config.New(path)
	if err != nil {
		log.Fatal(fmt.Errorf("error initializing config: %w", err))
	}

	l, err := logger.New(cfg.Logger.LogLevel)
	if err != nil {
		log.Fatal(fmt.Errorf("error initializing logger: %w", err))
	}

	ctx := context.Background()

	conn, err := postgresconn.New(ctx, cfg.Postgres)
	if err != nil {
		log.Fatal(fmt.Errorf("error initializing postgres connection: %w", err))
	}

	api, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(fmt.Errorf("error initializing botapi instance: %w", err))
	}

	es := emailSender.New(
		cfg.EmailSender.EmailToken,
		cfg.EmailSender.ClientEmail,
		cfg.EmailSender.HostEmail,
		cfg.EmailSender.SMTPAddress,
	)

	botCommander, err := commander.New(
		ctx,
		storage.New(conn),
		cfg.BotResponses,
		commander.WithLogger(l),
		commander.WithTgAPI(api),
		commander.WithEmailSender(es),
	)
	if err != nil {
		log.Fatal(fmt.Errorf("error initializing botcommander: %w", err))
	}

	var wg *sync.WaitGroup

	nctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	if err = botCommander.Start(ctx, wg); err != nil {
		log.Fatal(fmt.Errorf("error starting botcommander: %w", err))
	}

	ticker := time.NewTicker(time.Minute)

loop:
	for {
		select {
		case <-ticker.C:
			if err = botCommander.CheckForActualDate(ctx); err != nil {
				log.Fatal(err)
			}

			ticker.Reset(time.Minute)

		case <-nctx.Done():
			log.Println("graceful shutdown")

			break loop
		}
	}

	wg.Wait()
}
