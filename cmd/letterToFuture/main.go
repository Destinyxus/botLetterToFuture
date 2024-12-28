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

	botclient "github.com/Destinyxus/botLetterToFuture/internal/bot_client"
	commander "github.com/Destinyxus/botLetterToFuture/internal/bot_commander"
	"github.com/Destinyxus/botLetterToFuture/internal/config"
	emailclient "github.com/Destinyxus/botLetterToFuture/internal/email_client"
	"github.com/Destinyxus/botLetterToFuture/internal/logger"
	"github.com/Destinyxus/botLetterToFuture/internal/storage"
	"github.com/Destinyxus/botLetterToFuture/pkg/logruslog"
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

	l, err := logruslog.New(cfg.Logger.LogLevel)
	if err != nil {
		log.Fatal(fmt.Errorf("error initializing logrus: %w", err))
	}

	ctx := context.Background()

	conn, err := postgresconn.New(ctx, cfg.Postgres)
	if err != nil {
		log.Fatal(fmt.Errorf("error initializing postgres connection: %w", err))
	}

	botAPI, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(fmt.Errorf("error initializing botapi instance: %w", err))
	}

	botCommander, err := commander.New(
		ctx,
		storage.New(conn),
		cfg.BotResponses,
		commander.WithLogger(logger.New(l)),
		commander.WithTgAPI(botclient.New(botAPI)),
		commander.WithEmailClient(emailclient.New(
			cfg.EmailSender.EmailToken,
			cfg.EmailSender.ClientEmail,
			cfg.EmailSender.HostEmail,
			cfg.EmailSender.SMTPAddress,
		)),
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

	dateCheckInterval := time.NewTicker(cfg.DateCheckInterval)

loop:
	for {
		select {
		case <-dateCheckInterval.C:
			if err = botCommander.CheckForActualDate(ctx); err != nil {
				log.Fatal(fmt.Errorf("error checking for actual date: %w", err))
			}
		case <-nctx.Done():
			log.Println("graceful shutdown")

			break loop
		}
	}

	wg.Wait()
}
