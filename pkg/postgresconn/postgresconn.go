package postgresconn

import (
	"fmt"
	"github.com/Destinyxus/botLetterToFuture/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func New(cfg config.Config) (*sqlx.DB, error) {
	url := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
	)

	conn, err := sqlx.Connect("postgres", url)

	if err != nil {
		return &sqlx.DB{}, fmt.Errorf("initializing postgres connection: %w", err)
	}

	if err = conn.Ping(); err != nil {
		return &sqlx.DB{}, fmt.Errorf("pinging postgres db: %w", err)
	}

	migration := `CREATE TABLE if not exists letters
				(
				  id serial PRIMARY KEY,
   				  letter text not null,
    			  email varchar(255),
   				  date varchar,
    			  isActual bool
				)`

	_, err = conn.Exec(migration)
	if err != nil {
		return &sqlx.DB{}, err
	}

	return conn, nil
}
