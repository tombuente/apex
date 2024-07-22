package templates

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"strings"
)

var templateFuncs = template.FuncMap{
	"dict": func(values ...any) (map[string]any, error) {
		if len(values)%2 != 0 {
			return nil, errors.New("key value mismatch")
		}

		dict := make(map[string]any, len(values)/2)
		for i := 0; i < len(values); i += 2 {
			key, ok := values[i].(string)
			if !ok {
				return nil, errors.New("dict keys must be strings")
			}

			dict[key] = values[i+1]
		}
		return dict, nil
	},
}

func Load(templateFS fs.FS, service string) (map[string]*template.Template, error) {
	entries, err := fs.ReadDir(templateFS, fmt.Sprintf("templates/%v/views", service))
	if err != nil {
		return nil, err
	}

	var viewNames []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		viewNames = append(viewNames, entry.Name())
	}

	baseTemplate := "layout.html"
	layout := "templates/" + baseTemplate

	viewsPatternPrefix := fmt.Sprintf("templates/%v/views", service)
	componentsPattern := fmt.Sprintf("templates/%v/components/*.html", service)

	parsedTmpls := make(map[string]*template.Template)
	for _, viewName := range viewNames {
		tmplName := viewName[:strings.LastIndex(viewName, ".")] // stip file extension (.html) for template name

		// When creating a new template with template.New it should have the name of the default template. It gets executed with t.Execute.
		parsedTmpls[tmplName], err = template.New(baseTemplate).Funcs(templateFuncs).ParseFS(templateFS, layout, fmt.Sprintf("%v/%v", viewsPatternPrefix, viewName), componentsPattern)
		if err != nil {
			return nil, err
		}
	}

	slog.Info("Parsed templates", "service", service, "amount", len(parsedTmpls))
	return parsedTmpls, nil
}
