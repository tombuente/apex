package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/tombuente/apex/internal/xerrors"
)

func Exec(ctx context.Context, db *sqlx.DB, query string, args ...any) error {
	if _, err := db.ExecContext(context.Background(), query, args...); err != nil {
		return xerrors.ErrInternal
	}

	return nil
}

func One[T any](ctx context.Context, db *pgx.Conn, query string, args ...any) (T, error) {
	var i T
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return i, wrapError(err)
	}

	i, err = pgx.CollectOneRow[T](rows, pgx.RowToStructByNameLax)
	if err != nil {
		return i, wrapError(err)
	}

	return i, nil
}

func Many[T any](ctx context.Context, db *pgx.Conn, query string, args ...any) ([]T, error) {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, wrapError(err)
	}
	defer rows.Close()

	is, err := pgx.CollectRows[T](rows, pgx.RowToStructByNameLax)
	if err != nil {
		return nil, wrapError(err)
	}

	if len(is) == 0 {
		return nil, xerrors.ErrNotFound
	}

	return is, nil
}

func wrapError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return xerrors.Join(xerrors.ErrNotFound, err)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return fmt.Errorf("%w: postgres error: message: %v, code : %v", xerrors.ErrInternal, pgErr.Message, pgErr.Code)
	}

	return xerrors.Join(xerrors.ErrInternal, err)
}
