package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/tombuente/apex/internal/xerrors"
)

func One[T any](ctx context.Context, db *pgx.Conn, query string, args ...any) (T, error) {
	var defaultT T

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return defaultT, wrapError(err)
	}

	i, err := pgx.CollectOneRow[T](rows, pgx.RowToStructByNameLax)
	if err != nil {
		return defaultT, wrapError(err)
	}

	return i, nil
}

func Many[T any](ctx context.Context, db *pgx.Conn, query string, args ...any) ([]T, error) {
	var defaultT []T

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return defaultT, wrapError(err)
	}
	defer rows.Close()

	is, err := pgx.CollectRows[T](rows, pgx.RowToStructByNameLax)
	if err != nil {
		return defaultT, wrapError(err)
	}

	if len(is) == 0 {
		return defaultT, xerrors.ErrNotFound
	}

	return is, nil
}

func wrapError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return xerrors.Join(xerrors.ErrNotFound, err)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return fmt.Errorf("%w: postgres error: %v, code : %v", xerrors.ErrInternal, pgErr.Message, pgErr.Code)
	}

	return xerrors.Join(xerrors.ErrInternal, err)
}
