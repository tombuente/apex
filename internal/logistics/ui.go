package logistics

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"github.com/tombuente/apex/internal/templates"
	"github.com/tombuente/apex/internal/xerrors"
	"github.com/tombuente/apex/internal/xui"
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
	ui.templates, err = templates.Load("logistics")
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	r.Get("/", ui.indexView)

	r.Route("/items", func(r chi.Router) {
		r.Get("/{id}", xui.DetailViewWithData(ui.service.item, ui.newItemData, ui.templates["item-detail"]))
		r.Post("/{id}", xui.Update("/logistics/items", ui.service.updateItem))
		r.Get("/", xui.ListView(ui.newItemFilter, ui.service.items, ui.templates["item-list"]))
		r.Get("/new", xui.CreateViewWithData(ui.newItemData, ui.templates["item-create"]))
		r.Post("/", xui.Create("/logistics/items", ui.service.createItem))
	})

	r.Route("/plants", func(r chi.Router) {
		r.Get("/{id}", xui.DetailViewWithData(ui.service.plant, ui.newPlantData, ui.templates["plant-detail"]))
		r.Post("/{id}", xui.Update("/logistics/plants", ui.service.updatePlant))
		r.Get("/", xui.ListView(ui.newPlantFilter, ui.service.plants, ui.templates["plant-list"]))
		r.Get("/new", xui.CreateViewWithData(ui.newPlantData, ui.templates["plant-create"]))
		r.Post("/", xui.Create("/logistics/plants", ui.service.createPlant))
	})

	r.Route("/addresses", func(r chi.Router) {
		r.Get("/{id}", xui.DetailView(ui.service.address, ui.templates["address-detail"]))
		r.Post("/{id}", xui.Update("/logistics/addresses", ui.service.updateAddress))
		r.Get("/", xui.ListView(ui.newAddressFilter, ui.service.addresses, ui.templates["address-list"]))
		r.Get("/new", xui.CreateView[Address](ui.templates["address-create"]))
		r.Post("/", xui.Create("/logistics/addresses", ui.service.createAddress))
	})

	return r, nil
}

func (ui UI) indexView(w http.ResponseWriter, r *http.Request) {
	err := ui.templates["dashboard"].Execute(w, nil)
	if err != nil {
		fmt.Println(err)
	}
}

type itemData struct {
	Resource   *Item
	Categories []ItemCategory
}

func (ui UI) newItemData(ctx context.Context, item *Item) (itemData, error) {
	categories, err := ui.service.itemCategories(ctx)
	if err != nil {
		return itemData{}, err
	}

	return itemData{
		Resource:   item,
		Categories: categories,
	}, nil
}

func (ui UI) newItemFilter(ctx context.Context, values url.Values) (ItemFilter, error) {
	name := values.Get("name")
	sku := values.Get("sku")
	categoryID := values.Get("category_id")
	grossPrice := values.Get("gross_price")
	netPrice := values.Get("net_price")

	filter := ItemFilter{}

	if name != "" {
		filter.name = sql.NullString{Valid: true, String: name}
	}
	if sku != "" {
		filter.sku = sql.NullString{Valid: true, String: sku}
	}
	if categoryID != "" {
		categoryID, err := strconv.ParseInt(categoryID, 10, 64)
		if err != nil {
			return ItemFilter{}, fmt.Errorf("%w: unable to convert category id to integer", xerrors.ErrBadRequest)
		}
		filter.categoryID = sql.NullInt64{Valid: true, Int64: categoryID}
	}
	if grossPrice != "" {
		grossPrice, err := strconv.ParseInt(grossPrice, 10, 64)
		if err != nil {
			return ItemFilter{}, fmt.Errorf("%w: unable to convert gross price to integer", xerrors.ErrBadRequest)
		}
		filter.categoryID = sql.NullInt64{Valid: true, Int64: grossPrice}
	}
	if netPrice != "" {
		netPrice, err := strconv.ParseInt(netPrice, 10, 64)
		if err != nil {
			return ItemFilter{}, fmt.Errorf("%w: unable to convert net price to integer", xerrors.ErrBadRequest)
		}
		filter.categoryID = sql.NullInt64{Valid: true, Int64: netPrice}
	}

	return filter, nil
}

func (ui UI) newAddressFilter(ctx context.Context, values url.Values) (AddressFilter, error) {
	filter := AddressFilter{
		zip:       sql.NullString{},
		city:      sql.NullString{},
		street:    sql.NullString{},
		country:   sql.NullString{},
		latitude:  sql.NullFloat64{},
		longitude: sql.NullFloat64{},
	}

	return filter, nil
}

type plantData struct {
	Resource  *Plant
	Addresses []Address
}

func (ui UI) newPlantData(ctx context.Context, plant *Plant) (plantData, error) {
	addresses, err := ui.service.addresses(ctx, AddressFilter{})
	if err != nil {
		return plantData{}, err
	}

	return plantData{
		Resource:  plant,
		Addresses: addresses,
	}, nil
}

func (ui UI) newPlantFilter(ctx context.Context, values url.Values) (PlantFilter, error) {
	return PlantFilter{}, nil
}
