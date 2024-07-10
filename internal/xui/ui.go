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
	"github.com/tombuente/apex/internal/flash"
	"github.com/tombuente/apex/internal/xerrors"
)

var Decoder = schema.NewDecoder()

type Resource interface {
	GetID() string
	Redirect() string
}

type dataOne[R Resource] struct {
	Message  flash.Message
	Resource *R
}

type dataMany[R Resource] struct {
	Message   flash.Message
	Resources []R
}

func newDataOne[R Resource](ctx context.Context, w http.ResponseWriter, r *http.Request, resource *R) (dataOne[R], error) {
	return dataOne[R]{
		Message:  flash.Get(w, r),
		Resource: resource,
	}, nil
}

func newDataMany[R Resource](ctx context.Context, w http.ResponseWriter, r *http.Request, resources []R) (dataMany[R], error) {
	return dataMany[R]{
		Message:   flash.Get(w, r),
		Resources: resources,
	}, nil
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

		data, err := newDataMany(r.Context(), w, r, resources)
		if err != nil {
			slog.Error("Unable to create data", "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
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
	return DetailViewWithData(queryFunc, newDataOne, template)
}

func DetailViewWithData[R Resource, D any](
	queryFunc func(ctx context.Context, id int64) (R, error),
	makeDataFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request, resource *R) (D, error),
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

		data, err := makeDataFunc(r.Context(), w, r, &resource)
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
	return CreateViewWithData[R](newDataOne, template)
}

func CreateViewWithData[R Resource, D any](
	makeDataFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request, resource *R) (D, error),
	template *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := makeDataFunc(r.Context(), w, r, nil)
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

		flash.EntryCreated(w)
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

		flash.EntryUpdated(w)
		http.Redirect(w, r, item.Redirect(), http.StatusFound)
	}
}
