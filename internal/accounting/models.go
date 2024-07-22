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
	Description string `json:"description" form:"description"`
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

// Should only be embedded
type DocumentHeader struct {
	ID          int64  `json:"id" db:"id"`
	Description string `json:"description" db:"description"`
	Date        string `json:"date" db:"date"`
	PostingDate string `json:"posting_date" db:"posting_date"`
	Reference   string `json:"reference" db:"reference"`
	CurrencyID  int64  `json:"currency_id" db:"currency_id"`
}

// Should only be embedded
type DocumentPosition struct {
	ID          int64  `json:"id" db:"id"`
	DocumentID  int64  `json:"document_id" db:"document_id"`
	Description string `json:"description" db:"description"`
	AccountID   int64  `json:"account_id" db:"account_id"`
	TypeID      int64  `json:"type" db:"type_id"`
	Amount      int64  `json:"amount" db:"amount"`
}

type Document struct {
	DocumentHeader
	Positions []DocumentPosition `json:"positions"`
}

type DocumentParams struct {
	DocumentHeader
	Positions []DocumentPosition `json:"positions" form:"-"`
}

func (account Account) GetID() string {
	return strconv.FormatInt(account.ID, 10)
}

func (account Account) Redirect() string {
	return "/accounting/accounts" + account.GetID()
}

func (document Document) GetID() string {
	return strconv.FormatInt(document.ID, 10)
}

func (document Document) Redirect() string {
	return "/accounting/documents" + document.GetID()
}
