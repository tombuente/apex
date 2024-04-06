package templates

import (
	"fmt"
	"html/template"
	"reflect"
)

var funcs template.FuncMap

func init() {
	funcs = template.FuncMap{
		"has": has,
	}
}

func Load(templates map[string][]string) (map[string]*template.Template, error) {
	compiled := make(map[string]*template.Template)

	for name, files := range templates {
		var actualFiles []string
		for _, filename := range files {
			actualFiles = append(actualFiles, fmt.Sprintf("templates/%v.html", filename))
		}

		templ, err := template.ParseFiles(actualFiles...)
		if err != nil {
			return make(map[string]*template.Template), err
		}

		compiled[name] = templ
	}

	return compiled, nil
}

func has(value interface{}, name string) bool {
	reflectValue := reflect.ValueOf(value)
	if reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}

	if reflectValue.Kind() != reflect.Struct {
		return false
	}

	return reflectValue.FieldByName(name).IsValid()
}
