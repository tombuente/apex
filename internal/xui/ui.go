package xui

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"github.com/tombuente/apex/internal/xerrors"
)

var Decoder = schema.NewDecoder()

type Resource interface {
	IDString() string
}

type data[R Resource] struct {
	Resource *R
}

func newData[R Resource](ctx context.Context, resource *R) (data[R], error) {
	data := data[R]{
		Resource: resource,
	}

	return data, nil
}

func BasicView(template *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := template.Execute(w, nil)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
	}
}

func ListView[R Resource, F any](
	filterFunc func(ctx context.Context, values url.Values) (F, error),
	queryFunc func(ctx context.Context, filter F) ([]R, error),
	template *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		filter, err := filterFunc(r.Context(), r.URL.Query())

		resources, err := queryFunc(r.Context(), filter)
		if err != nil && !errors.Is(err, xerrors.ErrNotFound) {
			slog.Error("Unable to query resources", "error", err)
			xerrors.RenderHTML(w, "generic", err)
			return
		}

		ctx := struct {
			Resources []R
		}{
			Resources: resources,
		}

		fmt.Println()

		template.Execute(w, ctx)
		if err != nil {
			slog.Error("Unable to execute template", "error", err)
		}
	}
}

func DetailView[R Resource](
	queryFunc func(ctx context.Context, id int64) (R, error),
	template *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return DetailViewWithData(queryFunc, newData, template)
}

func DetailViewWithData[R Resource, D any](
	queryFunc func(ctx context.Context, id int64) (R, error),
	dataFunc func(ctx context.Context, resource *R) (D, error),
	template *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			http.Error(w, "malformatted id", http.StatusBadRequest)
			return
		}

		resource, err := queryFunc(r.Context(), id)
		if err != nil {
			xerrors.RenderHTML(w, "generic", err)
			return
		}

		data, err := dataFunc(r.Context(), &resource)
		if err != nil {
			xerrors.RenderHTML(w, "generic", err)
			return
		}

		err = template.Execute(w, data)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
	}
}

func CreateView[R Resource](template *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return CreateViewWithData[R](newData, template)
}

func CreateViewWithData[R Resource, D any](
	dataFunc func(ctx context.Context, resource *R) (D, error),
	template *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := dataFunc(r.Context(), nil)
		if err != nil {
			slog.Error("Unable to construct data object for create view", "error", err)
			xerrors.RenderHTML(w, "generic", err)
			return
		}

		err = template.Execute(w, data)
		if err != nil {
			slog.Error("Unable to execute template", "error", err)
		}
	}
}

func Create[R Resource, P any](
	redirect string,
	createFunc func(ctx context.Context, params P) (R, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "bad request", http.StatusInternalServerError)
			return
		}

		var params P
		err = Decoder.Decode(&params, r.PostForm)
		if err != nil {
			slog.Error("Unable to decode form", "error", err)
			http.Error(w, "unable to decode form", http.StatusBadRequest)
			return
		}

		item, err := createFunc(r.Context(), params)
		if err != nil {
			slog.Error("Unable to create entry in database", "error", err)
			xerrors.RenderHTML(w, "generic", err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("%v/%v", redirect, item.IDString()), http.StatusFound)
	}
}

func Update[R Resource, P any](
	redirect string,
	updateFunc func(ctx context.Context, id int64, params P) (R, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			http.Error(w, "malformatted id", http.StatusBadRequest)
			return
		}

		err = r.ParseForm()
		if err != nil {
			http.Error(w, "bad form", http.StatusBadRequest)
			return
		}

		var params P
		err = Decoder.Decode(&params, r.PostForm)
		if err != nil {
			http.Error(w, "unable to decode form", http.StatusBadRequest)
			return
		}

		updated, err := updateFunc(r.Context(), id, params)
		if err != nil {
			slog.Error("Unable to update", "error", err)
			xerrors.RenderHTML(w, "generic", err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("%v/%v", redirect, updated.IDString()), http.StatusFound)
	}
}
