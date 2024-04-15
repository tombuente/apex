package accounting

import (
	"context"
)

const name = "accounting"

type Service struct {
	db Database
}

func NewService(db Database) Service {
	return Service{
		db: db,
	}
}

func (s Service) account(ctx context.Context, id int64) (Account, error) {
	return s.db.account(ctx, id)
}

func (s Service) accounts(ctx context.Context, filter AccountFilter) ([]Account, error) {
	return s.db.accounts(ctx, filter)
}

func (s Service) createAccount(ctx context.Context, params AccountParams) (Account, error) {
	return s.db.createAccount(ctx, params)
}

func (s Service) updateAccount(ctx context.Context, id int64, params AccountParams) (Account, error) {
	return s.db.updateAccount(ctx, id, params)
}

// func (s Service) deleteAccount(ctx context.Context, id int64) error {
// 	return s.db.deleteAccount(ctx, id)
// }

func (s Service) currencies(ctx context.Context) ([]Currency, error) {
	return s.db.currencies(ctx)
}

func (s Service) documentPositionTypes(ctx context.Context) ([]DocumentPositionType, error) {
	return s.db.documentPositionTypes(ctx)
}

func (s Service) createDocument(ctx context.Context, params DocumentParams) ([]DocumentWithPositions, error) {
	return []DocumentWithPositions{}, nil
}
