package logistics

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/tombuente/apex/internal/xerr"
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
	if filter.gross_price != nil {
		builder = builder.Where(squirrel.Like{"gross_price": filter.gross_price})
	}
	if filter.net_price != nil {
		builder = builder.Where(squirrel.Like{"net_price": filter.net_price})
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

	row := db.db.QueryRowx(query, args...)

	var i Item
	err = row.StructScan(&i)
	if err != nil {
		return Item{}, xerr.Join(xerr.ErrInternal, err)
	}

	return i, nil
}
