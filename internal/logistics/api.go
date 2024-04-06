package logistics

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type API struct {
	service Service
}

func NewAPIRouter(service Service) chi.Router {
	a := API{
		service: service,
	}

	r := chi.NewRouter()
	r.Route("/items", func(r chi.Router) {
		r.Get("/{id}", a.getItem)
		r.Get("/", a.getItems)
		r.Post("/", a.postItem)
	})

	return r
}

func (a API) getItem(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	item, ok, err := a.service.item(r.Context(), ItemFilter{id: &id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	render.JSON(w, r, item)
}

func (a API) getItems(w http.ResponseWriter, r *http.Request) {
	items, ok, err := a.service.items(r.Context(), ItemFilter{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	render.JSON(w, r, items)
}

func (a API) postItem(w http.ResponseWriter, r *http.Request) {
	var params ItemParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	item, err := a.service.createItem(r.Context(), params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, item)
}
