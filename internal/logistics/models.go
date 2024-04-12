package logistics

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
	id         *int64
	name       *string
	sku        *string
	categoryID *int64
	GrossPrice *int64
	NetPrice   *int64
	limit      uint64
}

type Category struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type Address struct {
	ID           int64   `db:"id" json:"id"`
	ZIP          string  `db:"zip" json:"zip"`
	City         string  `db:"city" json:"city"`
	Street       string  `db:"street" json:"street"`
	StreetNumber string  `db:"street_number" json:"street_number"`
	Country      string  `db:"country" json:"country"`
	Longitude    float64 `db:"longitude" json:"longitude"`
	Latitude     float64 `db:"latitude" json:"latitude"`
}

type AddressParams struct {
	ZIP          string  `schema:"zip,required" json:"zip"`
	City         string  `schema:"city,required" json:"city"`
	Street       string  `schema:"street,required" json:"street"`
	StreetNumber string  `schema:"street_number,required" json:"street_number"`
	Country      string  `schema:"country,required" json:"country"`
	Longitude    float64 `schema:"longitude" json:"longitude"`
	Latitude     float64 `schema:"latitude" json:"latitude"`
}

type AddressFilter struct {
	id    *int64
	limit uint64
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
	id        *int64
	name      *string
	addressID *int64
	limit     uint64
}
