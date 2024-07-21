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

func (db Database) itemCategories(ctx context.Context) ([]ItemCategory, error) {
	const query = `
SELECT *
FROM logistics.item_categories
ORDER BY id ASC
`

	return database.Many[ItemCategory](ctx, db.db, query)
}

func (db Database) item(ctx context.Context, id int64) (Item, error) {
	const query = `
SELECT *
FROM logistics.items
WHERE id = $1
`

	return database.One[Item](ctx, db.db, query, id)
}

func (db Database) items(ctx context.Context, filter ItemFilter) ([]Item, error) {
	const query = `
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

	return database.Many[Item](ctx, db.db, query, filter.name, filter.sku, filter.categoryID, filter.grossPrice, filter.netPrice)
}

func (db Database) createItem(ctx context.Context, params ItemParams) (Item, error) {
	const query = `
INSERT INTO logistics.items (name, sku, category_id, gross_price, net_price)
VALUES ($1, $2, $3, $4, $5)
RETURNING *
`

	return database.One[Item](ctx, db.db, query, params.Name, params.SKU, params.CategoryID, params.GrossPrice, params.NetPrice)
}

func (db Database) updateItem(ctx context.Context, id int64, params ItemParams) (Item, error) {
	const query = `
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

	return database.One[Item](ctx, db.db, query, id, params.Name, params.SKU, params.CategoryID, params.GrossPrice, params.NetPrice)
}

func (db Database) address(ctx context.Context, id int64) (Address, error) {
	const query = `
SELECT *
FROM logistics.addresses
WHERE id = $1
`

	return database.One[Address](ctx, db.db, query, id)
}

func (db Database) addresses(ctx context.Context, filter AddressFilter) ([]Address, error) {
	const query = `
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

	return database.Many[Address](ctx, db.db, query, filter.zip, filter.city, filter.street, filter.country, filter.latitude, filter.longitude)
}

func (db Database) createAddress(ctx context.Context, params AddressParams) (Address, error) {
	const query = `
INSERT INTO logistics.addresses (zip, city, street, country, longitude, latitude)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *
`

	return database.One[Address](ctx, db.db, query, params.ZIP, params.City, params.Street, params.Country, params.Latitude, params.Longitude)
}

func (db Database) updateAddress(ctx context.Context, id int64, params AddressParams) (Address, error) {
	const query = `
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
	return database.One[Address](ctx, db.db, query, id, params.ZIP, params.City, params.Street, params.Country, params.Latitude, params.Longitude)
}

func (db Database) plant(ctx context.Context, id int64) (Plant, error) {
	const query = `
SELECT *
FROM logistics.plants
WHERE id = $1
`

	return database.One[Plant](ctx, db.db, query, id)
}

func (db Database) plants(ctx context.Context, filter PlantFilter) ([]Plant, error) {
	const query = `
SELECT *
FROM logistics.plants
WHERE
	(name       LIKE $1 OR $1 IS NULL) AND
	(address_id =    $2 OR $2 IS NULL)
ORDER BY id ASC
`

	return database.Many[Plant](ctx, db.db, query, filter.name, filter.addressID)
}

func (db Database) createPlant(ctx context.Context, params PlantParams) (Plant, error) {
	const query = `
INSERT INTO logistics.plants (name, address_id)
VALUES ($1, $2)
RETURNING *
`

	return database.One[Plant](ctx, db.db, query, params.Name, params.AddressID)
}

func (db Database) updatePlant(ctx context.Context, id int64, params PlantParams) (Plant, error) {
	const query = `
UPDATE logistics.plants
SET
	name       = $2,
	address_id = $3
WHERE id = $1
RETURNING *
`

	return database.One[Plant](ctx, db.db, query, id, params.Name, params.AddressID)
}
