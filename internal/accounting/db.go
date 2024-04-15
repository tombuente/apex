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

func NewDatabase(db *pgx.Conn) Database {
	return Database{
		db: db,
	}
}

const accountQuery = `
SELECT *
FROM accounting.accounts
WHERE id = $1
`

func (db Database) account(ctx context.Context, id int64) (Account, error) {
	return database.One[Account](ctx, db.db, accountQuery, id)
}

const accountsQuery = `
SELECT *
FROM accounting.accounts
WHERE (description LIKE $1 OR $1 IS NULL)
`

func (db Database) accounts(ctx context.Context, filter AccountFilter) ([]Account, error) {
	return database.Many[Account](ctx, db.db, accountsQuery, filter.Description)
}

const createAccountQuery = `
INSERT INTO accounting.accounts (description)
VALUES ($1)
RETURNING *
`

func (db Database) createAccount(ctx context.Context, params AccountParams) (Account, error) {
	return database.One[Account](ctx, db.db, createAccountQuery, params.Description)
}

const updateAccountQuery = `
UPDATE accounting.accounts
SET description = $2
WHERE id = $1
RETURNING *
`

func (db Database) updateAccount(ctx context.Context, id int64, params AccountParams) (Account, error) {
	return database.One[Account](ctx, db.db, updateAccountQuery, id, params.Description)
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

const currenciesQuery = `
SELECT *
FROM accounting.currencies
`

func (db Database) currencies(ctx context.Context) ([]Currency, error) {
	return database.Many[Currency](ctx, db.db, currenciesQuery)
}

const documentPositionTypesQuery = `
SELECT *
FROM accounting.document_position_types
`

func (db Database) documentPositionTypes(ctx context.Context) ([]DocumentPositionType, error) {
	return database.Many[DocumentPositionType](ctx, db.db, documentPositionTypesQuery)
}
