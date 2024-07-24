package accounting

import (
	"context"
	_ "embed"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/tombuente/apex/internal/database"
	"github.com/tombuente/apex/internal/xerrors"
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

func (db Database) document(ctx context.Context, id int64) (Document, error) {
	const headerQuery = `
SELECT *
FROM accounting.documents
WHERE id = $1
`

	header, err := database.One[DocumentHeader](ctx, db.db, headerQuery, id)
	if err != nil {
		return Document{}, err
	}

	const positionsQuery = `
SELECT *
FROM accounting.document_positions
WHERE document_id = $1
`

	positions, err := database.Many[DocumentPosition](ctx, db.db, positionsQuery, id)
	if err != nil && !errors.Is(err, xerrors.ErrNotFound) {
		return Document{}, err
	}

	return Document{DocumentHeader: header, Positions: positions}, nil
}

func (db Database) documents(ctx context.Context, _ DocumentFilter) ([]Document, error) {
	const query = `
SELECT *
FROM accounting.documents
`

	return database.Many[Document](ctx, db.db, query)
}

func (db Database) createDocument(ctx context.Context, params DocumentParams) (Document, error) {
	// TODO: use transaction to make sure documents are not created without positions and vise versa

	const documentHeaderQuery = `
INSERT INTO accounting.documents (date, posting_date, reference, description, currency_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *
`
	documentHeader, err := database.One[DocumentHeader](ctx, db.db, documentHeaderQuery, params.Date, params.PostingDate, params.Reference, params.Description, params.CurrencyID)
	if err != nil {
		return Document{}, err
	}

	const documentPositionsQuery = `
INSERT INTO accounting.document_positions (document_id, account_id, description, type_id, amount)
VALUES ($1, $2, $3, $4, $5)
RETURNING *
`
	var documentPositions []DocumentPosition

	for _, posParams := range params.Positions {
		documentPosition, err := database.One[DocumentPosition](ctx, db.db, documentPositionsQuery, documentHeader.ID, posParams.AccountID, posParams.Description, posParams.TypeID, posParams.Amount)
		if err != nil {
			return Document{}, err
		}

		documentPositions = append(documentPositions, documentPosition)
	}

	return Document{DocumentHeader: documentHeader, Positions: documentPositions}, nil
}
