package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/Destinyxus/botLetterToFuture/internal/bot_commander"
	commander "github.com/Destinyxus/botLetterToFuture/internal/bot_commander"
	"github.com/jackc/pgx/v5"
)

type storage struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) commander.Storage {
	return &storage{
		conn: conn,
	}
}

func (s *storage) InsertLetter(ctx context.Context, letter, email string, date time.Time) error {
	_, err := s.conn.Exec(ctx, `
		insert into letters
		(
		  letter,
		  email,
		  date,
		  isActual
		)
		values 
		(
		  @letter,
		  @email,
		  @date,
		  @isActual
		)`, pgx.NamedArgs{
		"letter":   letter,
		"email":    email,
		"date":     date,
		"isActual": true,
	})
	if err != nil {
		return fmt.Errorf("error inserting into letters: %w", err)
	}

	return nil
}

func (s *storage) GetLetter(ctx context.Context, date time.Time) ([]bot_commander.Letter, error) {
	var letter []bot_commander.Letter

	if err := s.conn.QueryRow(ctx, `
		select * 
		from letters 
		where date = @date 
		and isActual is true
	`, pgx.NamedArgs{"date": date}).Scan(letter); err != nil {
		return nil, fmt.Errorf("error selecting from letters: %w", err)
	}

	return letter, nil
}

func (s *storage) DeprecateLetter(ctx context.Context, date time.Time) error {
	if _, err := s.conn.Exec(ctx, `
		update letters 
		set isActual = false 
		where date = @date`, pgx.NamedArgs{"date": date}); err != nil {
		return fmt.Errorf("error deprecating letters: %w", err)
	}

	return nil
}

func (s *storage) GetActualDates(ctx context.Context) (map[time.Time]struct{}, error) {
	rows, err := s.conn.Query(ctx, `
		select distinct date 
		from letters 
		where date >= CURRENT_DATE`,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting actual dates: %w", err)
	}
	defer rows.Close()

	dates := []time.Time{}

	for rows.Next() {
		var date time.Time

		if err = rows.Scan(&date); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		dates = append(dates, date)
	}

	dateIndexes := make(map[time.Time]struct{})

	for _, d := range dates {
		dateIndexes[d] = struct{}{}
	}

	return dateIndexes, nil
}
