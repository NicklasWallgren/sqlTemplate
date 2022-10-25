package pkg_test

import (
	"embed"
	"testing"
	"text/template"

	"github.com/NicklasWallgren/sqlTemplate/pkg"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/*.tsql
var fs embed.FS // nolint: varnamelen

func TestNewQueryTemplateEngine(t *testing.T) {
	pkg.NewQueryTemplateEngine()
}

func TestNewQueryTemplateEngine_WithTemplateFunctions(t *testing.T) {
	if sqlT := pkg.NewQueryTemplateEngine(pkg.WithTemplateFunctions(template.FuncMap{})); sqlT == nil {
		t.Fatal()
	}
}

func TestQueryTemplateEngine_Register(t *testing.T) {
	sqlT := pkg.NewQueryTemplateEngine(pkg.WithTemplateFunctions(template.FuncMap{}))

	if err := sqlT.Register("users", fs, ".tsql"); err != nil {
		t.Fatal()
	}
}

func TestQueryTemplateEngine_Parse(t *testing.T) {
	sqlT := pkg.NewQueryTemplateEngine(pkg.WithTemplateFunctions(template.FuncMap{}))

	if err := sqlT.Register("users", fs, ".tsql"); err != nil {
		t.Fatal()
	}

	template, err := sqlT.Parse("users", "findUsers")
	assert.Nil(t, err)
	assert.Equal(t, "\n    SELECT *\n    FROM users\n", template.GetQuery())
	assert.Equal(t, []any{}, template.GetParams())
}

func TestQueryTemplateEngine_ParseWithValuesFromMap(t *testing.T) {
	sqlT := pkg.NewQueryTemplateEngine(pkg.WithTemplateFunctions(template.FuncMap{}))

	if err := sqlT.Register("users", fs, ".tsql"); err != nil {
		t.Fatal()
	}

	criteria := map[string]interface{}{"ID": 1, "Order": "id"}

	template, err := sqlT.ParseWithValuesFromMap("users", "findById", criteria)
	assert.Nil(t, err)
	assert.Equal(t, "\n    SELECT *\n    FROM users\n    WHERE id=?\n    ORDER BY id\n", template.GetQuery())
	assert.Equal(t, []any{1}, template.GetParams())
}

func TestQueryTemplateEngine_ParseWithValuesFromStruct(t *testing.T) {
	sqlT := pkg.NewQueryTemplateEngine(pkg.WithTemplateFunctions(template.FuncMap{}))

	if err := sqlT.Register("users", fs, ".tsql"); err != nil {
		t.Fatal()
	}

	type searchCriteria struct { // nolint:govet
		ID    int
		Order string
	}

	criteria := searchCriteria{ID: 1, Order: "id"}

	template, err := sqlT.ParseWithValuesFromStruct("users", "findById", criteria)
	assert.Nil(t, err)
	assert.Equal(t, "\n    SELECT *\n    FROM users\n    WHERE id=?\n    ORDER BY id\n", template.GetQuery())
	assert.Equal(t, []any{1}, template.GetParams())
}
