package accounting

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tombuente/apex/internal/flash"
	"github.com/tombuente/apex/internal/templates"
	"github.com/tombuente/apex/internal/xerrors"
	"github.com/tombuente/apex/internal/xui"
)

type UI struct {
	service   Service
	templates map[string]*template.Template
}

type accountsData struct {
	Message   flash.Message
	Resources []Account
}

type documentData struct {
	Message       flash.Message
	Resource      *Document
	Accounts      []Account
	Currencies    []Currency
	PositionTypes []DocumentPositionType
	Positions     []DocumentPosition
}

func NewUIRouter(templateFS fs.FS, service Service) (*chi.Mux, error) {
	ui := UI{
		service:   service,
		templates: make(map[string]*template.Template),
	}

	var err error
	ui.templates, err = templates.Load(templateFS, "accounting")
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()

	r.Get("/", xui.BasicView(ui.templates["dashboard"]))

	r.Route("/accounts", func(r chi.Router) {
		r.Get("/{id}", xui.Detail(ui.service.account, ui.templates["account-detail"]))
		r.Post("/{id}", xui.Update(ui.service.updateAccount))
		r.Get("/", ui.accountListView)
		r.Get("/new", xui.CreateView[Account](ui.templates["account-create"]))
		r.Post("/", xui.Create(ui.service.createAccount))
	})

	r.Route("/documents", func(r chi.Router) {
		r.Get("/{id}", xui.DetailWithAdditionalData(ui.service.document, ui.additionalDocumentData, ui.templates["document-detail"]))
		r.Get("/new", xui.CreateViewWithData(ui.additionalDocumentData, ui.templates["document-create"]))
		r.Get("/", xui.ListView(ui.makeDocumentFilter, ui.service.documents, ui.templates["document-list"]))
		r.Post("/", xui.CreateWithFormParser(parseDocumentForm, ui.service.createDocument))
		// r.Post("/verify", ui.vertifyDocumentViewHTMX)
	})

	return r, nil
}

func (ui UI) accountListView(w http.ResponseWriter, r *http.Request) {
	accounts, err := ui.service.accounts(r.Context(), AccountFilter{})
	if err != nil && !errors.Is(err, xerrors.ErrNotFound) {
		slog.Error("Unable to get accounts from database", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := accountsData{
		Message:   flash.Get(w, r),
		Resources: accounts,
	}

	err = ui.templates["account-list"].Execute(w, data)
	if err != nil {
		slog.Error("Unable to execute template", "error", err)
	}
}

func (ui UI) additionalDocumentData(ctx context.Context, w http.ResponseWriter, r *http.Request, document *Document) (documentData, error) {
	accounts, err := ui.service.accounts(ctx, AccountFilter{})
	if err != nil {
		return documentData{}, nil
	}

	currencies, err := ui.service.currencies(ctx)
	if err != nil {
		return documentData{}, err
	}

	documentPositionTypes, err := ui.service.documentPositionTypes(ctx)
	if err != nil {
		return documentData{}, err
	}

	return documentData{
		Message:       flash.Get(w, r),
		Resource:      document,
		Accounts:      accounts,
		Currencies:    currencies,
		PositionTypes: documentPositionTypes,
	}, nil
}

func (ui UI) makeDocumentFilter(ctx context.Context, values url.Values) (DocumentFilter, error) {
	// TODO: Remove dummy filter
	return DocumentFilter{}, nil
}

func parseDocumentForm(values url.Values) (DocumentParams, error) {
	currencyID, err := strconv.ParseInt(values.Get("currency_id"), 10, 64)
	if err != nil {
		return DocumentParams{}, fmt.Errorf("unable to parse document currency_id to integer: %w", err)
	}

	header := DocumentHeaderParams{
		Description: values.Get("description"),
		Date:        values.Get("date"),
		PostingDate: values.Get("posting_date"),
		Reference:   values.Get("reference"),
		CurrencyID:  currencyID,
	}

	var positions []DocumentPositionParams
	for i := 0; i < len(values["positions[].account_id"]); i++ {
		description := values["positions[].description"][i]
		accountID, err := strconv.ParseInt(values["positions[].account_id"][i], 10, 64)
		if err != nil {
			return DocumentParams{}, fmt.Errorf("unable to parse position account_id to integer: %w", err)
		}
		typeID, err := strconv.ParseInt(values["positions[].type_id"][i], 10, 64)
		if err != nil {
			return DocumentParams{}, fmt.Errorf("unable to parse position type_id to integer: %w", err)
		}
		amount, err := strconv.ParseInt(values["positions[].amount"][i], 10, 64)
		if err != nil {
			return DocumentParams{}, fmt.Errorf("unable to parse position amount to integer: %w", err)
		}

		positions = append(positions, DocumentPositionParams{Description: description, AccountID: accountID, TypeID: typeID, Amount: amount})
	}

	return DocumentParams{DocumentHeaderParams: header, Positions: positions}, nil
}

// func (ui UI) vertifyDocumentViewHTMX(w http.ResponseWriter, r *http.Request) {
// 	err := r.ParseForm()
// 	if err != nil {
// 		http.Error(w, "bad request", http.StatusInternalServerError)
// 		return
// 	}

// 	if err := ui.templates["document-health.htmx"].ExecuteTemplate(w, "health", nil); err != nil {
// 		slog.Error("Unable to execute template", "error", err)
// 		http.Error(w, "internal error", http.StatusInternalServerError)
// 		return
// 	}
// }
