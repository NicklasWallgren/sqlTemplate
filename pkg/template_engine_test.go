package pkg

import (
	"fmt"
	"os"
	"testing"
	"text/template"
)

func TestNewQueryTemplateEngine(t *testing.T) {
	NewQueryTemplateEngine()
}

func TestNewQueryTemplateEngine_WithTemplateFunctions(t *testing.T) {
	_ = NewQueryTemplateEngine(WithTemplateFunctions(template.FuncMap{}))
}

func TestQueryTemplateEngine_Register(t *testing.T) {

	sqlT := NewQueryTemplateEngine()

	fs := os.DirFS("/Users/nicklaswallgren/Dropbox/projects/golang/sql-named-parameters/examples/map_as_param/queries")
	
	if err := sqlT.Register("users", fs, ".tsql"); err != nil {
		panic(err)
	}

	criteria := map[string]interface{}{"Id": "1", "Order": "id"}

	tmpl, err := sqlT.ParseWithValuesFromMap("users", "findById", criteria)
	if err != nil {
		panic(err)
	}

	// nolint:forbidigo
	fmt.Printf("query %v\n", tmpl.GetQuery())
	// nolint:forbidigo
	fmt.Printf("query parameters %v\n", tmpl.GetParams())

}

func TestQueryTemplateEngine_Parse(t *testing.T) {

}

func TestQueryTemplateEngine_ParseWithValuesFromMap(t *testing.T) {

}

func TestQueryTemplateEngine_ParseWithValuesFromStruct(t *testing.T) {

}
