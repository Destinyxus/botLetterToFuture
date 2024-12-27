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
	url := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
	)

	pgxCfg, err := pgx.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("pgx.parseconfig error: %w", err)
	}

	conn, err := pgx.ConnectConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("pgx.connect error: %w", err)
	}

	file, err := os.ReadFile(cfg.MigrationPath)
	if err != nil {
		return nil, fmt.Errorf("createtable error: %w", err)
	}

	_, err = conn.Exec(ctx, string(file))
	if err != nil {
		return nil, fmt.Errorf("conn.exec, creating table error: %w", err)
	}

	return conn, nil
}
