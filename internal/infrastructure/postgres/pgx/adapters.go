package pgx

import (
	"errors"
	"fmt"

	"github.com/Newo123/payment-bot-go/internal/infrastructure/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type pgxRows struct {
	pgx.Rows
}

type pgxRow struct {
	pgx.Row
}

func (r pgxRow) Scan(dest ...any) error {
	if err := r.Row.Scan(dest...); err != nil {
		return mapErrors(err)
	}

	return nil
}

type pgxCommandTag struct {
	pgconn.CommandTag
}

func mapErrors(err error) error {
	const (
		pgxViolatesForeignKeyErrorCode = "23503"
		pgxViolatesUniqueKeyErrorCode  = "23505"
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return postgres.ErrNoRows
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgxViolatesForeignKeyErrorCode:
			return fmt.Errorf("%v: %w", err, postgres.ErrViolatesForeignKey)
		case pgxViolatesUniqueKeyErrorCode:
			return fmt.Errorf("%v: %w", err, postgres.ErrViolatesUniqueKey)
		}
	}

	return fmt.Errorf("%v: %w", err, postgres.ErrUnknown)
}
