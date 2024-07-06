package xui

import (
	"context"
	"errors"
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
	GetID() string
	Redirect() string
}

type data[R Resource] struct {
	Resource *R
}

type dataMany[R Resource] struct {
	Resources *[]R
}

func newData[R Resource](ctx context.Context, resource *R) (data[R], error) {
	data := data[R]{
		Resource: resource,
	}

	return data, nil
}

func BasicView(template *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := template.Execute(w, nil); err != nil {
			slog.Error("Unable to render template", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}
}

func ListView[R Resource, F any](
	makeFilterFunc func(ctx context.Context, values url.Values) (F, error),
	queryFunc func(ctx context.Context, filter F) ([]R, error),
	template *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resourceFilter, err := makeFilterFunc(r.Context(), r.URL.Query())
		if err != nil {
			slog.Error("Unable to create filter with makeFilterFunc", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		resources, err := queryFunc(r.Context(), resourceFilter)
		if err != nil && !errors.Is(err, xerrors.ErrNotFound) {
			slog.Error("Unable to query resources", "error", err)
			http.Error(w, "unable to query resources", http.StatusInternalServerError)
			return
		}

		var data dataMany[R]
		if !errors.Is(err, xerrors.ErrNotFound) {
			data.Resources = &resources
		}

		if err := template.Execute(w, data); err != nil {
			slog.Error("Unable to execute template", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
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
	makeDataFunc func(ctx context.Context, resource *R) (D, error),
	template *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			http.Error(w, "malformatted id", http.StatusBadRequest)
			return
		}

		resource, err := queryFunc(r.Context(), id)
		if err != nil {
			slog.Error("Unable to query resource", "error", err)
			msg, code := xerrors.HttpInfo(err)
			http.Error(w, msg, code)
			return
		}

		data, err := makeDataFunc(r.Context(), &resource)
		if err != nil {
			slog.Error("Unable to make data", "error", err)
			msg, code := xerrors.HttpInfo(err)
			http.Error(w, msg, code)
			return
		}

		if err := template.Execute(w, data); err != nil {
			slog.Error("Unable to execute template", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}
}

func CreateView[R Resource](template *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return CreateViewWithData[R](newData, template)
}

func CreateViewWithData[R Resource, D any](
	makeDataFunc func(ctx context.Context, resource *R) (D, error),
	template *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := makeDataFunc(r.Context(), nil)
		if err != nil {
			slog.Error("Unable to make data", "error", err)
			msg, err := xerrors.HttpInfo(err)
			http.Error(w, msg, err)
			return
		}

		if err := template.Execute(w, data); err != nil {
			slog.Error("Unable to execute template", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}
}

func Create[R Resource, P any](
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
			http.Error(w, "unable to create resource", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, item.Redirect(), http.StatusFound)
	}
}

func Update[R Resource, P any](
	updateFunc func(ctx context.Context, id int64, params P) (R, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			http.Error(w, "malformatted id", http.StatusBadRequest)
			return
		}

		err = r.ParseForm()
		if err != nil {
			http.Error(w, "unable to parse form", http.StatusBadRequest)
			return
		}

		var params P
		err = Decoder.Decode(&params, r.PostForm)
		if err != nil {
			slog.Error("Unable to ", "error", err)
			http.Error(w, "unable to decode form", http.StatusBadRequest)
			return
		}

		item, err := updateFunc(r.Context(), id, params)
		if err != nil {
			slog.Error("Unable to update resources", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, item.Redirect(), http.StatusFound)
	}
}
