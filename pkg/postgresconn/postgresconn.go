package postgresconn

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
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

	migration := `CREATE TABLE if not exists letters
				(
				  id serial PRIMARY KEY,
   				  letter text not null,
    			  email varchar(255),
   				  date varchar,
    			  isActual bool
				)`

	_, err = conn.Exec(ctx, migration)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
