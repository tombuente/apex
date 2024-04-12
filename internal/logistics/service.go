package logistics

import (
	"context"
)

const name = "logistics"

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

func (s Service) updateItem(ctx context.Context, id int64, params ItemParams) (Item, error) {
	return s.db.updateItem(ctx, id, params)
}

func (s Service) address(ctx context.Context, filter AddressFilter) (Address, bool, error) {
	return s.db.address(ctx, filter)
}

func (s Service) addresses(ctx context.Context, filter AddressFilter) ([]Address, bool, error) {
	return s.db.addresses(ctx, filter)
}

func (s Service) createAddress(ctx context.Context, params AddressParams) (Address, error) {
	return s.db.createAddress(ctx, params)
}

func (s Service) updateAddress(ctx context.Context, id int64, params AddressParams) (Address, error) {
	return s.db.updateAddress(ctx, id, params)
}

func (s Service) plant(ctx context.Context, filter PlantFilter) (Plant, bool, error) {
	return s.db.plant(ctx, filter)
}

func (s Service) plants(ctx context.Context, filter PlantFilter) ([]Plant, bool, error) {
	return s.db.plants(ctx, filter)
}

func (s Service) createPlant(ctx context.Context, params PlantParams) (Plant, error) {
	return s.db.createPlants(ctx, params)
}

func (s Service) updatePlant(ctx context.Context, id int64, params PlantParams) (Plant, error) {
	return s.db.updatePlant(ctx, id, params)
}
