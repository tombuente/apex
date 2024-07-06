package logistics

import (
	"context"
)

type Service struct {
	db Database
}

func MakeService(db Database) Service {
	return Service{
		db: db,
	}
}

func (s Service) itemCategories(ctx context.Context) ([]ItemCategory, error) {
	return s.db.itemCategories(ctx)
}

func (s Service) item(ctx context.Context, id int64) (Item, error) {
	return s.db.item(ctx, id)
}

func (s Service) items(ctx context.Context, filter ItemFilter) ([]Item, error) {
	return s.db.items(ctx, filter)
}

func (s Service) createItem(ctx context.Context, params ItemParams) (Item, error) {
	return s.db.createItem(ctx, params)
}

func (s Service) updateItem(ctx context.Context, id int64, params ItemParams) (Item, error) {
	return s.db.updateItem(ctx, id, params)
}

func (s Service) address(ctx context.Context, id int64) (Address, error) {
	return s.db.address(ctx, id)
}

func (s Service) addresses(ctx context.Context, filter AddressFilter) ([]Address, error) {
	return s.db.addresses(ctx, filter)
}

func (s Service) createAddress(ctx context.Context, params AddressParams) (Address, error) {
	return s.db.createAddress(ctx, params)
}

func (s Service) updateAddress(ctx context.Context, id int64, params AddressParams) (Address, error) {
	return s.db.updateAddress(ctx, id, params)
}

func (s Service) plant(ctx context.Context, id int64) (Plant, error) {
	return s.db.plant(ctx, id)
}

func (s Service) plants(ctx context.Context, filter PlantFilter) ([]Plant, error) {
	return s.db.plants(ctx, filter)
}

func (s Service) createPlant(ctx context.Context, params PlantParams) (Plant, error) {
	return s.db.createPlant(ctx, params)
}

func (s Service) updatePlant(ctx context.Context, id int64, params PlantParams) (Plant, error) {
	return s.db.updatePlant(ctx, id, params)
}
