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

	templs["address-list"] = []string{"layout", "logistics/views/address-list"}
	templs["address-detail"] = []string{"layout", "logistics/views/address-detail", "logistics/components/address-form"}
	templs["address-create"] = []string{"layout", "logistics/views/address-create", "logistics/components/address-form"}

	templs["plant-list"] = []string{"layout", "logistics/views/plant-list"}
	templs["plant-detail"] = []string{"layout", "logistics/views/plant-detail", "logistics/components/plant-form"}
	templs["plant-create"] = []string{"layout", "logistics/views/plant-create", "logistics/components/plant-form"}

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
	r.Route("/addresses", func(r chi.Router) {
		r.Get("/{id}", ui.addressDetailView)
		r.Post("/{id}", ui.addressUpdate)
		r.Get("/", ui.addressListView)
		r.Post("/", ui.addressCreate)
		r.Get("/new", ui.addressCreateView)
	})
	r.Route("/plants", func(r chi.Router) {
		r.Get("/{id}", ui.plantDetailView)
		r.Post("/{id}", ui.plantUpdate)
		r.Get("/", ui.plantListView)
		r.Post("/", ui.plantCreate)
		r.Get("/new", ui.plantCreateView)
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
		filter.GrossPrice = &grossPrice
	}
	if hasNetPrice {
		filter.NetPrice = &netPrice
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

	_, err = ui.service.updateItem(r.Context(), id, params)
	if err != nil {
		slog.Info("Unable to update database with params", "endpoint", "itemUpdate", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/logistics/items", http.StatusFound)
}

func (ui UI) addressListView(w http.ResponseWriter, r *http.Request) {
	var filter AddressFilter

	addresses, _, err := ui.service.addresses(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx := struct {
		Addresses []Address
	}{
		Addresses: addresses,
	}

	err = ui.templates["address-list"].Execute(w, ctx)
	if err != nil {
		slog.Error("Unable to execute template", "error", err)
	}
}

func (ui UI) addressDetailView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	address, ok, err := ui.service.address(r.Context(), AddressFilter{id: &id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ctx := struct {
		Action  string
		Address *Address
	}{
		Action:  fmt.Sprintf("/logistics/addresses/%v", id),
		Address: &address,
	}

	err = ui.templates["address-detail"].Execute(w, ctx)
	if err != nil {
		slog.Error("Unable to execute template", "error", err)
	}
}

func (ui UI) addressCreateView(w http.ResponseWriter, r *http.Request) {
	ctx := struct {
		Action  string
		Address *Address
	}{
		Action:  "/logistics/addresses",
		Address: nil,
	}

	err := ui.templates["address-create"].Execute(w, ctx)
	if err != nil {
		slog.Error("Unable to execute template", "error", err)
	}
}

func (ui UI) addressCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var params AddressParams
	err = decoder.Decode(&params, r.PostForm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = ui.service.createAddress(r.Context(), params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/logistics/addresses", http.StatusFound)
}

func (ui UI) addressUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		slog.Info("Unable to parse id", "endpoint", "addressUpdate", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var params AddressParams
	err = decoder.Decode(&params, r.PostForm)
	if err != nil {
		slog.Info("Unable to decode form", "endpoint", "addressUpdate", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = ui.service.updateAddress(r.Context(), id, params)
	if err != nil {
		slog.Info("Unable to update database with params", "endpoint", "addressUpdate", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/logistics/addresses", http.StatusFound)
}

func (ui UI) plantListView(w http.ResponseWriter, r *http.Request) {
	var filter PlantFilter

	plants, _, err := ui.service.plants(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx := struct {
		Plants []Plant
	}{
		Plants: plants,
	}

	err = ui.templates["plant-list"].Execute(w, ctx)
	if err != nil {
		slog.Error("Unable to execute template", "error", err)
	}
}

func (ui UI) plantDetailView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	plant, ok, err := ui.service.plant(r.Context(), PlantFilter{id: &id})
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
		Plant  *Plant
	}{
		Action: fmt.Sprintf("/logistics/plants/%v", id),
		Plant:  &plant,
	}

	err = ui.templates["plant-detail"].Execute(w, ctx)
	if err != nil {
		slog.Error("Unable to execute template", "error", err)
	}
}

func (ui UI) plantCreateView(w http.ResponseWriter, r *http.Request) {
	ctx := struct {
		Action string
		Plant  *Address
	}{
		Action: "/logistics/plants",
		Plant:  nil,
	}

	err := ui.templates["plant-create"].Execute(w, ctx)
	if err != nil {
		slog.Error("Unable to execute template", "error", err)
	}
}

func (ui UI) plantCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var params PlantParams
	err = decoder.Decode(&params, r.PostForm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = ui.service.createPlant(r.Context(), params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/logistics/plants", http.StatusFound)
}

func (ui UI) plantUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		slog.Info("Unable to parse id", "endpoint", "addressUpdate", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var params PlantParams
	err = decoder.Decode(&params, r.PostForm)
	if err != nil {
		slog.Info("Unable to decode form", "endpoint", "addressUpdate", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = ui.service.db.updatePlant(r.Context(), id, params)
	if err != nil {
		slog.Info("Unable to update database with params", "endpoint", "addressUpdate", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/logistics/plants", http.StatusFound)
}
