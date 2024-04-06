package logistics

import (
	"context"
)

type Service struct {
	db Database
}

func NewService(db Database) Service {
	return Service{
		db: db,
	}
}

func (s Service) item(ctx context.Context, filter ItemFilter) (Item, bool, error) {
	return s.db.item(ctx, filter)
}

func (s Service) items(ctx context.Context, filter ItemFilter) ([]Item, bool, error) {
	return s.db.items(ctx, filter)
}

func (s Service) createItem(ctx context.Context, params ItemParams) (Item, error) {
	return s.db.createItem(ctx, params)
}
