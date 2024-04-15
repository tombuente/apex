package accounting

import (
	"database/sql"
	"strconv"
)

type Account struct {
	ID          int64  `json:"id" db:"id"`
	Description string `json:"description" db:"description"`
}

type AccountParams struct {
	Description string `json:"description" schema:"description"`
}

type AccountFilter struct {
	Description sql.NullString
}

type Currency struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	ISO  string `json:"iso" db:"iso"`
}

type DocumentPositionType struct {
	ID          int64  `json:"id" db:"id"`
	Description string `json:"description" db:"description"`
}

type Document struct {
	ID          int64  `json:"id" db:"id"`
	Description string `json:"description" db:"description"`
	Date        string `json:"date" db:"date"`
	PostingDate string `json:"posting_date" db:"posting_date"`
	Reference   string `json:"reference" db:"reference"`
	CurrencyID  int64  `json:"currency_id" db:"currency_id"`
}

func (document Document) IDString() string {
	return strconv.FormatInt(document.ID, 10)
}

type DocumentPosition struct {
	ID          int64  `json:"id" db:"id"`
	DocumentID  int64  `json:"document_id" db:"document_id"`
	Description string `json:"description" db:"description"`
	AccountID   int64  `json:"account_id" db:"account_id"`
	Type        int64  `json:"type" db:"type"`
	Amount      int64  `json:"amount" db:"amount"`
}

func (account Account) IDString() string {
	return strconv.FormatInt(account.ID, 10)
}
