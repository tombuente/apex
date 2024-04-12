package accounting

import (
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"github.com/tombuente/apex/internal/templates"
	"github.com/tombuente/apex/internal/xerrors"
)

var decoder = schema.NewDecoder()

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

	r.Get("/", ui.indexView)

	r.Route("/accounts", func(r chi.Router) {
		r.Get("/{id}", ui.accountDetailView)
		// r.Post("/{id}", ui.itemUpdate)
		r.Get("/", ui.accountListView)
		// r.Post("/", ui.itemCreate)
		// r.Get("/new", ui.itemCreateView)
	})

	return r, nil
}

func (ui UI) indexView(w http.ResponseWriter, r *http.Request) {
	err := ui.templates["dashboard"].Execute(w, nil)
	if err != nil {
		fmt.Println(err)
	}
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

func (ui UI) accountDetailView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	account, err := ui.service.account(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx := struct {
		Action  string
		Account *Account
	}{
		Action:  fmt.Sprintf("/accounting/accounts/%v", id),
		Account: &account,
	}

	err = ui.templates["account-detail"].Execute(w, ctx)
	if err != nil {
		slog.Error("Unable to execute template", "error", err)
	}
}
