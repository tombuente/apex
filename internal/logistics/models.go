package logistics

import (
	"database/sql"
	"strconv"
)

type Item struct {
	ID         int64  `db:"id" json:"id"`
	Name       string `db:"name" json:"name"`
	SKU        string `db:"sku" json:"sku"`
	CategoryID int64  `db:"category_id" json:"category_id"`
	GrossPrice int64  `db:"gross_price" json:"gross_price"`
	NetPrice   int64  `db:"net_price" json:"net_price"`
}

type ItemParams struct {
	Name       string `json:"name" schema:"name,required"`
	SKU        string `json:"sku" schema:"sku,required"`
	CategoryID int64  `json:"category_id" schema:"category_id,required"`
	GrossPrice int64  `json:"gross_price" schema:"gross_price"`
	NetPrice   int64  `json:"net_price" schema:"net_price"`
}

type ItemFilter struct {
	id         sql.NullInt64
	name       sql.NullString
	sku        sql.NullString
	categoryID sql.NullInt64
	grossPrice sql.NullInt64
	netPrice   sql.NullInt64
}

type ItemCategory struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type Address struct {
	ID        int64   `db:"id" json:"id"`
	ZIP       string  `db:"zip" json:"zip"`
	City      string  `db:"city" json:"city"`
	Street    string  `db:"street" json:"street"`
	Country   string  `db:"country" json:"country"`
	Latitude  float64 `db:"latitude" json:"latitude"`
	Longitude float64 `db:"longitude" json:"longitude"`
}

type AddressParams struct {
	ZIP       string  `schema:"zip,required" json:"zip"`
	City      string  `schema:"city,required" json:"city"`
	Street    string  `schema:"street,required" json:"street"`
	Country   string  `schema:"country,required" json:"country"`
	Latitude  float64 `schema:"latitude" json:"latitude"`
	Longitude float64 `schema:"longitude" json:"longitude"`
}

type AddressFilter struct {
	zip       sql.NullString
	city      sql.NullString
	street    sql.NullString
	country   sql.NullString
	longitude sql.NullFloat64
	latitude  sql.NullFloat64
}

type Plant struct {
	ID        int64  `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	AddressID int64  `db:"address_id" json:"address_id"`
}

type PlantParams struct {
	Name      string `schema:"name" json:"name"`
	AddressID int64  `schema:"address_id, required" json:"address_id"`
}

type PlantFilter struct {
	name      sql.NullString
	addressID sql.NullInt64
}

func (item Item) GetID() string {
	return strconv.FormatInt(item.ID, 10)
}

func (item Item) Redirect() string {
	return "/logistics/items/" + item.GetID()
}

func (itemCategory ItemCategory) GetID() string {
	return strconv.FormatInt(itemCategory.ID, 10)
}

func (itemCategory ItemCategory) Redirect() string {
	return "/logistics/itemcategories/" + itemCategory.GetID()
}

func (address Address) GetID() string {
	return strconv.FormatInt(address.ID, 10)
}

func (address Address) Redirect() string {
	return "/logistics/addresses/" + address.GetID()
}

func (plant Plant) GetID() string {
	return strconv.FormatInt(plant.ID, 10)
}

func (plant Plant) Redirect() string {
	return "/logistics/plants/" + plant.GetID()
}
