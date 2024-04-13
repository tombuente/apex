package xui

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"github.com/tombuente/apex/internal/xerrors"
)

var decoder = schema.NewDecoder()

type Resource interface {
	IDString() string
}

func BasicView(template *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := template.Execute(w, nil)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
	}
}

func DetailView[R Resource](queryFunc func(ctx context.Context, id int64) (R, error), template *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			http.Error(w, "malformatted id", http.StatusBadRequest)
			return
		}

		items, err := queryFunc(r.Context(), id)
		if err != nil {
			xerrors.RenderHTML(w, "generic", err)
			return
		}

		data := struct {
			Resource *R
		}{
			Resource: &items,
		}

		err = template.Execute(w, data)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
	}
}

func CreateView[R Resource](template *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Resource *R
		}{}

		err := template.Execute(w, data)
		if err != nil {
			slog.Error("Unable to execute template", "error", err)
		}
	}
}

func Create[R Resource, P any](redirect string, createFunc func(ctx context.Context, params P) (R, error)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "bad request", http.StatusInternalServerError)
			return
		}

		var params P
		err = decoder.Decode(&params, r.PostForm)
		if err != nil {
			http.Error(w, "unable to decode form", http.StatusBadRequest)
			return
		}

		item, err := createFunc(r.Context(), params)
		if err != nil {
			xerrors.RenderHTML(w, "generic", err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("%v/%v", redirect, item.IDString()), http.StatusFound)
	}
}

func Update[R Resource, P any](redirect string, updateFunc func(ctx context.Context, id int64, params P) (R, error)) func(w http.ResponseWriter, r *http.Request) {
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
		err = decoder.Decode(&params, r.PostForm)
		if err != nil {
			http.Error(w, "unable to decode form", http.StatusBadRequest)
			return
		}

		updated, err := updateFunc(r.Context(), id, params)
		if err != nil {
			xerrors.RenderHTML(w, "generic", err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("%v/%v", redirect, updated.IDString()), http.StatusFound)
	}
}
