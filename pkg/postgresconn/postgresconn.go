package postgresconn

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

func New(ctx context.Context, addr string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, addr)

	if err != nil {
		return nil, fmt.Errorf("initializing postgres connection: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if err = conn.Ping(pingCtx); err != nil {
		return nil, fmt.Errorf("pinging postgres db: %w", err)
	}

	return conn, nil
}
