package logistics

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/tombuente/apex/internal/xerr"
)

var (
	itemsTable     = fmt.Sprintf("%v_items", appName)
	addressesTable = fmt.Sprintf("%v_addresses", appName)
	plantsTable    = fmt.Sprintf("%v_plants", appName)
)

type Database struct {
	db *sqlx.DB
}

func NewDatabase(db *sqlx.DB) Database {
	return Database{
		db: db,
	}
}

func (db Database) item(ctx context.Context, filter ItemFilter) (Item, bool, error) {
	filter.limit = 1

	items, ok, err := db.items(ctx, filter)
	if err != nil {
		return Item{}, ok, nil
	}
	if !ok {
		return Item{}, ok, nil
	}

	return items[0], ok, nil
}

func (db Database) items(ctx context.Context, filter ItemFilter) ([]Item, bool, error) {
	builder := squirrel.Select("*").From(itemsTable)

	if filter.id != nil {
		builder = builder.Where(squirrel.Like{"id": filter.id})
	}
	if filter.name != nil {
		builder = builder.Where(squirrel.Like{"name": filter.name})
	}
	if filter.sku != nil {
		builder = builder.Where(squirrel.Like{"sku": filter.sku})
	}
	if filter.categoryID != nil {
		builder = builder.Where(squirrel.Like{"category_id": filter.categoryID})
	}
	if filter.GrossPrice != nil {
		builder = builder.Where(squirrel.Like{"gross_price": filter.GrossPrice})
	}
	if filter.NetPrice != nil {
		builder = builder.Where(squirrel.Like{"net_price": filter.NetPrice})
	}

	if filter.limit != 0 {
		builder = builder.Limit(filter.limit)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return []Item{}, false, xerr.Join(xerr.ErrInternal, err)
	}

	rows, err := db.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return []Item{}, false, xerr.Join(xerr.ErrInternal, err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var i Item
		err := rows.StructScan(&i)
		if err != nil {
			return []Item{}, false, xerr.Join(xerr.ErrInternal, err)
		}

		items = append(items, i)
	}

	if len(items) == 0 {
		return []Item{}, false, nil
	}

	return items, true, nil
}

func (db Database) createItem(ctx context.Context, params ItemParams) (Item, error) {
	query, args, err := squirrel.Insert(itemsTable).
		Columns("name, sku, category_id, gross_price, net_price").
		Values(params.Name, params.SKU, params.CategoryID, params.GrossPrice, params.NetPrice).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return Item{}, xerr.Join(xerr.ErrInternal, err)
	}

	row := db.db.QueryRowxContext(ctx, query, args...)

	var i Item
	err = row.StructScan(&i)
	if err != nil {
		return Item{}, xerr.Join(xerr.ErrInternal, err)
	}

	return i, nil
}

func (db Database) updateItem(ctx context.Context, id int64, params ItemParams) (Item, error) {
	query, args, err := squirrel.Update(itemsTable).
		Set("name", params.Name).
		Set("sku", params.SKU).
		Set("category_id", params.CategoryID).
		Set("gross_price", params.GrossPrice).
		Set("net_price", params.NetPrice).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING *").
		ToSql()

	if err != nil {
		return Item{}, xerr.Join(xerr.ErrInternal, err)
	}

	row := db.db.QueryRowxContext(ctx, query, args...)

	var i Item
	err = row.StructScan(&i)
	if err != nil {
		return Item{}, xerr.Join(xerr.ErrInternal, err)
	}

	return i, nil
}

func (db Database) address(ctx context.Context, filter AddressFilter) (Address, bool, error) {
	filter.limit = 1

	address, ok, err := db.addresses(ctx, filter)
	if err != nil {
		return Address{}, ok, nil
	}
	if !ok {
		return Address{}, ok, nil
	}

	return address[0], ok, nil
}

func (db Database) addresses(ctx context.Context, filter AddressFilter) ([]Address, bool, error) {
	builder := squirrel.Select("*").From(addressesTable)

	if filter.id != nil {
		builder = builder.Where(squirrel.Like{"id": filter.id})
	}

	if filter.limit != 0 {
		builder = builder.Limit(filter.limit)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return []Address{}, false, xerr.Join(xerr.ErrInternal, err)
	}

	rows, err := db.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return []Address{}, false, xerr.Join(xerr.ErrInternal, err)
	}
	defer rows.Close()

	var addresses []Address
	for rows.Next() {
		var i Address
		err := rows.StructScan(&i)
		if err != nil {
			return []Address{}, false, xerr.Join(xerr.ErrInternal, err)
		}

		addresses = append(addresses, i)
	}

	if len(addresses) == 0 {
		return []Address{}, false, nil
	}

	return addresses, true, nil
}

func (db Database) createAddress(ctx context.Context, params AddressParams) (Address, error) {
	query, args, err := squirrel.Insert(addressesTable).
		Columns("zip, city, street, street_number, country, longitude, latitude").
		Values(params.ZIP, params.City, params.Street, params.StreetNumber, params.Country, params.Latitude, params.Longitude).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return Address{}, xerr.Join(xerr.ErrInternal, err)
	}

	row := db.db.QueryRowxContext(ctx, query, args...)

	var i Address
	err = row.StructScan(&i)
	if err != nil {
		return Address{}, xerr.Join(xerr.ErrInternal, err)
	}

	return i, nil
}

func (db Database) updateAddress(ctx context.Context, id int64, params AddressParams) (Address, error) {
	query, args, err := squirrel.Update(addressesTable).
		Set("zip", params.ZIP).
		Set("city", params.City).
		Set("street", params.Street).
		Set("street_number", params.StreetNumber).
		Set("country", params.Country).
		Set("longitude", params.Longitude).
		Set("latitude", params.Latitude).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING *").
		ToSql()

	if err != nil {
		return Address{}, xerr.Join(xerr.ErrInternal, err)
	}

	row := db.db.QueryRowxContext(ctx, query, args...)

	var i Address
	err = row.StructScan(&i)
	if err != nil {
		return Address{}, xerr.Join(xerr.ErrInternal, err)
	}

	return i, nil
}

func (db Database) plant(ctx context.Context, filter PlantFilter) (Plant, bool, error) {
	filter.limit = 1

	plants, ok, err := db.plants(ctx, filter)
	if err != nil {
		return Plant{}, ok, nil
	}
	if !ok {
		return Plant{}, ok, nil
	}

	return plants[0], ok, nil
}

func (db Database) plants(ctx context.Context, filter PlantFilter) ([]Plant, bool, error) {
	builder := squirrel.Select("*").From(plantsTable)

	if filter.id != nil {
		builder = builder.Where(squirrel.Like{"id": filter.id})
	}
	if filter.name != nil {
		builder = builder.Where(squirrel.Like{"name": filter.name})
	}
	if filter.addressID != nil {
		builder = builder.Where(squirrel.Like{"address_id": filter.addressID})
	}

	if filter.limit != 0 {
		builder = builder.Limit(filter.limit)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return []Plant{}, false, xerr.Join(xerr.ErrInternal, err)
	}

	rows, err := db.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return []Plant{}, false, xerr.Join(xerr.ErrInternal, err)
	}
	defer rows.Close()

	var plants []Plant
	for rows.Next() {
		var i Plant
		err := rows.StructScan(&i)
		if err != nil {
			return []Plant{}, false, xerr.Join(xerr.ErrInternal, err)
		}

		plants = append(plants, i)
	}

	if len(plants) == 0 {
		return []Plant{}, false, nil
	}

	return plants, true, nil
}

func (db Database) createPlants(ctx context.Context, params PlantParams) (Plant, error) {
	query, args, err := squirrel.Insert(plantsTable).
		Columns("name, address_id").
		Values(params.Name, params.AddressID).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return Plant{}, xerr.Join(xerr.ErrInternal, err)
	}

	row := db.db.QueryRowxContext(ctx, query, args...)

	var i Plant
	err = row.StructScan(&i)
	if err != nil {
		return Plant{}, xerr.Join(xerr.ErrInternal, err)
	}

	return i, nil
}

func (db Database) updatePlant(ctx context.Context, id int64, params PlantParams) (Plant, error) {
	query, args, err := squirrel.Update(plantsTable).
		Set("name", params.Name).
		Set("address_id", params.AddressID).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING *").
		ToSql()

	if err != nil {
		return Plant{}, xerr.Join(xerr.ErrInternal, err)
	}

	row := db.db.QueryRowxContext(ctx, query, args...)

	var i Plant
	err = row.StructScan(&i)
	if err != nil {
		return Plant{}, xerr.Join(xerr.ErrInternal, err)
	}

	return i, nil
}
