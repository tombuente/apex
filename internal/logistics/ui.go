package logistics

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"github.com/tombuente/apex/internal/templates"
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

	templs := make(map[string][]string)
	templs["dashboard"] = []string{"layout", "logistics/views/dashboard"}
	templs["item-list"] = []string{"layout", "logistics/views/item-list"}
	templs["item-detail"] = []string{"layout", "logistics/views/item-detail", "logistics/components/item-form"}
	templs["item-create"] = []string{"layout", "logistics/views/item-create", "logistics/components/item-form"}

	var err error
	ui.templates, err = templates.Load(templs)
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	r.Get("/", ui.indexView)
	r.Route("/items", func(r chi.Router) {
		r.Get("/{id}", ui.itemDetailView)
		r.Post("/{id}", ui.itemUpdate)
		r.Get("/", ui.itemListView)
		r.Post("/", ui.itemCreate)
		r.Get("/new", ui.itemCreateView)
	})

	return r, nil
}

func (ui UI) indexView(w http.ResponseWriter, r *http.Request) {
	err := ui.templates["dashboard"].Execute(w, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func (ui UI) itemListView(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	sku := r.URL.Query().Get("sku")
	hasCategoryID := true
	categoryID, err := strconv.ParseInt(r.URL.Query().Get("category_id"), 10, 64)
	if err != nil {
		hasCategoryID = false
	}
	hasGrossPrice := true
	grossPrice, err := strconv.ParseInt(r.URL.Query().Get("gross_price"), 10, 64)
	if err != nil {
		hasGrossPrice = false
	}
	hasNetPrice := true
	netPrice, err := strconv.ParseInt(r.URL.Query().Get("net_price"), 10, 64)
	if err != nil {
		hasNetPrice = false
	}

	var filter ItemFilter
	if name != "" {
		filter.name = &name
	}
	if sku != "" {
		filter.sku = &sku
	}
	if hasCategoryID {
		filter.categoryID = &categoryID
	}
	if hasGrossPrice {
		filter.gross_price = &grossPrice
	}
	if hasNetPrice {
		filter.net_price = &netPrice
	}

	items, _, err := ui.service.items(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx := struct {
		Items []Item
	}{
		Items: items,
	}

	err = ui.templates["item-list"].Execute(w, ctx)
	if err != nil {
		slog.Error("Unable to execute template", "error", err)
	}
}

func (ui UI) itemDetailView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	item, ok, err := ui.service.item(r.Context(), ItemFilter{id: &id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ctx := struct {
		Action string
		Item   *Item
	}{
		Action: fmt.Sprintf("/logistics/items/%v", id),
		Item:   &item,
	}

	err = ui.templates["item-detail"].Execute(w, ctx)
	if err != nil {
		slog.Error("Unable to execute template", "error", err)
	}
}

func (ui UI) itemCreateView(w http.ResponseWriter, r *http.Request) {
	ctx := struct {
		Action string
		Item   *Item
	}{
		Action: "/logistics/items",
		Item:   nil,
	}

	err := ui.templates["item-create"].Execute(w, ctx)
	if err != nil {
		slog.Error("Unable to execute template", "error", err)
	}
}

func (ui UI) itemCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var params ItemParams
	err = decoder.Decode(&params, r.PostForm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = ui.service.createItem(r.Context(), params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/logistics/items", http.StatusFound)
}

func (ui UI) itemUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		slog.Info("Unable to parse id", "endpoint", "itemUpdate", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var params ItemParams
	err = decoder.Decode(&params, r.PostForm)
	if err != nil {
		slog.Info("Unable to decode form", "endpoint", "itemUpdate", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = ui.service.db.updateItem(r.Context(), id, params)
	if err != nil {
		slog.Info("Unable to update database with params", "endpoint", "itemUpdate", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/logistics/items", http.StatusFound)
}
