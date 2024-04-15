package accounting

import (
	"context"
	"errors"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tombuente/apex/internal/templates"
	"github.com/tombuente/apex/internal/xerrors"
	"github.com/tombuente/apex/internal/xui"
)

type UI struct {
	service   Service
	templates map[string]*template.Template
}

func NewUIRouter(service Service) (chi.Router, error) {
	ui := UI{
		service:   service,
		templates: make(map[string]*template.Template),
	}

	var err error
	ui.templates, err = templates.Load("accounting")
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()

	r.Get("/", xui.BasicView(ui.templates["dashboard"]))

	r.Route("/accounts", func(r chi.Router) {
		r.Get("/{id}", xui.DetailView(ui.service.account, ui.templates["account-detail"]))
		r.Post("/{id}", xui.Update("/accounting/accounts", ui.service.updateAccount))
		r.Get("/", ui.accountListView)
		r.Get("/new", xui.CreateView[Account](ui.templates["account-create"]))
		r.Post("/", xui.Create("/accounting/accounts", ui.service.createAccount))
	})

	r.Route("/documents", func(r chi.Router) {
		r.Get("/new", xui.CreateViewWithData(ui.newDocumentData, ui.templates["document-create"]))
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

	ctx := struct {
		Accounts []Account
	}{
		Accounts: accounts,
	}

	err = ui.templates["account-list"].Execute(w, ctx)
	if err != nil {
		slog.Error("Unable to execute template", "error", err)
	}
}

type documentData struct {
	Resource      *Document
	Accounts      []Account
	Currencies    []Currency
	PositionTypes []DocumentPositionType
	Positions     []DocumentPosition
}

func (ui UI) newDocumentData(ctx context.Context, document *Document) (documentData, error) {
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
		Resource:      document,
		Accounts:      accounts,
		Currencies:    currencies,
		PositionTypes: documentPositionTypes,
		Positions:     nil,
	}, nil
}
