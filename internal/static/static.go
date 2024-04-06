package static

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func NewStaticRouter(root http.FileSystem) (chi.Router, error) {
	path := "/"

	if strings.ContainsAny(path, "{}*") {
		return nil, errors.New("FileServer does not permit any URL parameters")
	}

	r := chi.NewRouter()

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})

	return r, nil
}
