package accounting

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5"
	"github.com/tombuente/apex/internal/database"
)

//go:embed schema.sql
var Schema string

//go:embed fixture.sql
var Fixture string

type Database struct {
	db *pgx.Conn
}

func MakeDatabase(db *pgx.Conn) Database {
	return Database{
		db: db,
	}
}

func (db Database) account(ctx context.Context, id int64) (Account, error) {
	const query = `
SELECT *
FROM accounting.accounts
WHERE id = $1
`

	return database.One[Account](ctx, db.db, query, id)
}

func (db Database) accounts(ctx context.Context, filter AccountFilter) ([]Account, error) {
	const query = `
SELECT *
FROM accounting.accounts
WHERE (description LIKE $1 OR $1 IS NULL)
`

	return database.Many[Account](ctx, db.db, query, filter.Description)
}

func (db Database) createAccount(ctx context.Context, params AccountParams) (Account, error) {
	const query = `
INSERT INTO accounting.accounts (description)
VALUES ($1)
RETURNING *
`

	return database.One[Account](ctx, db.db, query, params.Description)
}

func (db Database) updateAccount(ctx context.Context, id int64, params AccountParams) (Account, error) {
	const query = `
UPDATE accounting.accounts
SET description = $2
WHERE id = $1
RETURNING *
`

	return database.One[Account](ctx, db.db, query, id, params.Description)
}

// const deleteAccountQuery = `
// UPDATE accounting.accounts
// SET description = $2
// WHERE id = $1
// RETURNING *
// `

// func (db Database) deleteAccount(ctx context.Context, id int64) error {
// 	return database.Exec(ctx, db.db, deleteAccountQuery, id)
// }

func (db Database) currencies(ctx context.Context) ([]Currency, error) {
	const query = `
SELECT *
FROM accounting.currencies
`

	return database.Many[Currency](ctx, db.db, query)
}

func (db Database) documentPositionTypes(ctx context.Context) ([]DocumentPositionType, error) {
	const query = `
SELECT *
FROM accounting.document_position_types
`

	return database.Many[DocumentPositionType](ctx, db.db, query)
}
