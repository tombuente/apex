package logistics

import "fmt"

const appName = "logistics"

var (
	itemsTable = fmt.Sprintf("%v_items", appName)
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
	id          *int64
	name        *string
	sku         *string
	categoryID  *int64
	gross_price *int64
	net_price   *int64
	limit       uint64
}

type Category struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}
