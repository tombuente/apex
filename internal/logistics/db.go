package logistics

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

const itemCategoriesQuery = `
SELECT *
FROM logistics.item_categories
ORDER BY id ASC
`

func (db Database) itemCategories(ctx context.Context) ([]ItemCategory, error) {
	return database.Many[ItemCategory](ctx, db.db, itemCategoriesQuery)
}

const itemQuery = `
SELECT *
FROM logistics.items
WHERE id = $1
`

func (db Database) item(ctx context.Context, id int64) (Item, error) {
	return database.One[Item](ctx, db.db, itemQuery, id)
}

const itemsQuery = `
SELECT *
FROM logistics.items
WHERE
	(name        LIKE $1 OR $1 IS NULL) AND
	(sku         LIKE $2 OR $2 IS NULL) AND
	(category_id = $3 OR $3 IS NULL) AND
	(gross_price = $4 OR $4 IS NULL) AND
	(net_price   = $5 OR $5 IS NULL)
ORDER BY id ASC
`

func (db Database) items(ctx context.Context, filter ItemFilter) ([]Item, error) {
	return database.Many[Item](ctx, db.db, itemsQuery, filter.name, filter.sku, filter.categoryID, filter.grossPrice, filter.netPrice)
}

const createItemQuery = `
INSERT INTO logistics.items (name, sku, category_id, gross_price, net_price)
VALUES ($1, $2, $3, $4, $5)
RETURNING *
`

func (db Database) createItem(ctx context.Context, params ItemParams) (Item, error) {
	return database.One[Item](ctx, db.db, createItemQuery, params.Name, params.SKU, params.CategoryID, params.GrossPrice, params.NetPrice)
}

const updateItemQuery = `
UPDATE logistics.items
SET
	name        = $2,
	sku         = $3,
	category_id = $4,
	gross_price = $5,
	net_price   = $6
WHERE id = $1
RETURNING *
`

func (db Database) updateItem(ctx context.Context, id int64, params ItemParams) (Item, error) {
	return database.One[Item](ctx, db.db, updateItemQuery, id, params.Name, params.SKU, params.CategoryID, params.GrossPrice, params.NetPrice)
}

const addressQuery = `
SELECT *
FROM logistics.addresses
WHERE id = $1
`

func (db Database) address(ctx context.Context, id int64) (Address, error) {
	return database.One[Address](ctx, db.db, addressQuery, id)
}

const addressesQuery = `
SELECT *
FROM logistics.addresses
WHERE
	(zip       LIKE $1 OR $1 IS NULL) AND
	(city      LIKE $2 OR $2 IS NULL) AND
	(street    LIKE $3 OR $3 IS NULL) AND
	(country   LIKE $4 OR $4 IS NULL) AND
	(latitude  =    $5 OR $5 IS NULL) AND
	(longitude =    $6 OR $6 IS NULL)
ORDER BY id ASC
`

func (db Database) addresses(ctx context.Context, filter AddressFilter) ([]Address, error) {
	return database.Many[Address](ctx, db.db, addressesQuery, filter.zip, filter.city, filter.street, filter.country, filter.latitude, filter.longitude)
}

const createAddressQuery = `
INSERT INTO logistics.addresses (zip, city, street, country, longitude, latitude)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *
`

func (db Database) createAddress(ctx context.Context, params AddressParams) (Address, error) {
	return database.One[Address](ctx, db.db, createAddressQuery, params.ZIP, params.City, params.Street, params.Country, params.Latitude, params.Longitude)
}

const updateAddressQuery = `
UPDATE logistics.addresses
SET
	zip       = $2,
	city      = $3,
	street    = $4,
	country   = $5,
	latitude  = $6,
	longitude = $7
WHERE id = $1
RETURNING *
`

func (db Database) updateAddress(ctx context.Context, id int64, params AddressParams) (Address, error) {
	return database.One[Address](ctx, db.db, updateAddressQuery, id, params.ZIP, params.City, params.Street, params.Country, params.Latitude, params.Longitude)
}

const plantQuery = `
SELECT *
FROM logistics.plants
WHERE id = $1
`

func (db Database) plant(ctx context.Context, id int64) (Plant, error) {
	return database.One[Plant](ctx, db.db, plantQuery, id)
}

const plantsQuery = `
SELECT *
FROM logistics.plants
WHERE
	(name       LIKE $1 OR $1 IS NULL) AND
	(address_id =    $2 OR $2 IS NULL)
ORDER BY id ASC
`

func (db Database) plants(ctx context.Context, filter PlantFilter) ([]Plant, error) {
	return database.Many[Plant](ctx, db.db, plantsQuery, filter.name, filter.addressID)
}

const createPlantQuery = `
INSERT INTO logistics.plants (name, address_id)
VALUES ($1, $2)
RETURNING *
`

func (db Database) createPlant(ctx context.Context, params PlantParams) (Plant, error) {
	return database.One[Plant](ctx, db.db, createPlantQuery, params.Name, params.AddressID)
}

const updatePlantQuery = `
UPDATE logistics.plants
SET
	name       = $2,
	address_id = $3
WHERE id = $1
RETURNING *
`

func (db Database) updatePlant(ctx context.Context, id int64, params PlantParams) (Plant, error) {
	return database.One[Plant](ctx, db.db, updatePlantQuery, id, params.Name, params.AddressID)
}
