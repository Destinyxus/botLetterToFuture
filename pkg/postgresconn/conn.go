package postgresconn

import (
	"context"
	"fmt"
	"os"

	"github.com/Destinyxus/botLetterToFuture/internal/config"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

func New(ctx context.Context, cfg config.Postgres) (*pgx.Conn, error) {
	pgxCfg, err := pgx.ParseConfig(fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
	))
	if err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	conn, err := pgx.ConnectConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("error connecting: %w", err)
	}

	file, err := os.ReadFile(cfg.MigrationPath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	_, err = conn.Exec(ctx, string(file))
	if err != nil {
		return nil, fmt.Errorf("error creating table: %w", err)
	}

	return conn, nil
}
