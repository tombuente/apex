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
